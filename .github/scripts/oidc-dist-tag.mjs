#!/usr/bin/env node
// Promote an npm dist-tag using GitHub Actions OIDC → npm registry token exchange.
//
// `npm dist-tag` cannot authenticate via Trusted Publishing OIDC directly
// (npm/cli#8547). This mirrors the Nuxt release workaround: exchange the
// Actions id-token for a short-lived npm token, then PUT the dist-tags API.
//
// Usage:
//   node .github/scripts/oidc-dist-tag.mjs <package> <version> <tag>
//   node .github/scripts/oidc-dist-tag.mjs @sparrow-community/wasm 2026.9.30-alpha.1 latest

const [packageName, version, tag] = process.argv.slice(2);

if (!packageName || !version || !tag) {
  console.error("usage: oidc-dist-tag.mjs <package> <version> <tag>");
  process.exit(2);
}

async function githubOidcIdToken() {
  const requestUrl = process.env.ACTIONS_ID_TOKEN_REQUEST_URL;
  const requestToken = process.env.ACTIONS_ID_TOKEN_REQUEST_TOKEN;
  if (!requestUrl || !requestToken) {
    throw new Error(
      "ACTIONS_ID_TOKEN_REQUEST_URL/TOKEN unset — run only on GitHub Actions with id-token: write",
    );
  }
  const url = `${requestUrl}&audience=${encodeURIComponent("npm:registry.npmjs.org")}`;
  const res = await fetch(url, {
    headers: { authorization: `Bearer ${requestToken}` },
  });
  if (!res.ok) {
    throw new Error(`GitHub OIDC id-token failed: ${res.status} ${res.statusText}`);
  }
  const body = await res.json();
  if (!body?.value) {
    throw new Error("GitHub OIDC response missing value");
  }
  return body.value;
}

async function exchangeNpmToken(idToken, pkg) {
  const encoded = pkg.replace("/", "%2f");
  const url = `https://registry.npmjs.org/-/npm/v1/oidc/token/exchange/package/${encoded}`;
  const res = await fetch(url, {
    method: "POST",
    headers: { authorization: `Bearer ${idToken}` },
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(
      `npm OIDC exchange failed for ${pkg}: ${res.status} ${res.statusText} — ${text}`,
    );
  }
  const body = await res.json();
  if (!body?.token) {
    throw new Error("npm OIDC exchange response missing token");
  }
  return body.token;
}

async function putDistTag(npmToken, pkg, ver, distTag) {
  const encoded = pkg.replace("/", "%2f");
  const url = `https://registry.npmjs.org/-/package/${encoded}/dist-tags/${encodeURIComponent(distTag)}`;
  const res = await fetch(url, {
    method: "PUT",
    headers: {
      authorization: `Bearer ${npmToken}`,
      "content-type": "application/json",
    },
    body: JSON.stringify(ver),
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(
      `dist-tag PUT ${distTag} → ${pkg}@${ver} failed: ${res.status} ${res.statusText} — ${text}`,
    );
  }
}

const idToken = await githubOidcIdToken();
const npmToken = await exchangeNpmToken(idToken, packageName);
await putDistTag(npmToken, packageName, version, tag);
console.log(`dist-tag ${tag} → ${packageName}@${version}`);
