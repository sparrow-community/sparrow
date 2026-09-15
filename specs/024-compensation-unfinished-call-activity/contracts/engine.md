# Engine contract

## Runtime

- Parent compensate throw may publish `compensate_unfinished_child` for each unfinished Call Activity under throw scope (or the activityRef target).
- Child: terminate active tokens; queue CompensationSubs; run handlers; then terminate process and ResumeParent(Completed=false).
- Parent: WaitingChildren > 0 blocks advance; each host terminate after child done decrements; then advanceCompensation / Complete throw.

## Deploy

No new Deploy acceptance rules.
