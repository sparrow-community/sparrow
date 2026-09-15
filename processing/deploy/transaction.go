package deploy

import (
	"fmt"
	"strings"

	"github.com/sparrow-community/sparrow/bpmn/element"
	eventv1 "github.com/sparrow-community/sparrow/protocol/gen/go/event/v1"
)

const transactionMethodCompensate = "##Compensate"

// childScopes returns FlowElements of direct embedded SubProcesses and Transactions.
func childScopes(fe *element.FlowElements) []*element.FlowElements {
	if fe == nil {
		return nil
	}
	out := make([]*element.FlowElements, 0, len(fe.SubProcesses)+len(fe.Transactions))
	for i := range fe.SubProcesses {
		out = append(out, &fe.SubProcesses[i].FlowElements)
	}
	for i := range fe.Transactions {
		out = append(out, &fe.Transactions[i].FlowElements)
	}
	return out
}

func validateTransaction(tx element.Transaction, insideTransaction bool) error {
	if insideTransaction {
		return fmt.Errorf("UNSUPPORTED_ELEMENT: nested transaction %q is not supported", tx.ID)
	}
	method := strings.TrimSpace(tx.Method)
	if method != "" && method != transactionMethodCompensate {
		return fmt.Errorf("UNSUPPORTED_ELEMENT: transaction %q method %q not supported (only ##Compensate)", tx.ID, method)
	}
	if len(tx.LoopCharacteristicsElements.MultielementLoopCharacteristics) > 0 ||
		len(tx.LoopCharacteristicsElements.StandardLoopCharacteristics) > 0 {
		return fmt.Errorf("UNSUPPORTED_ELEMENT: multi-instance / standard loop on transaction %q is not supported", tx.ID)
	}
	if len(tx.StartEvents) == 0 {
		return fmt.Errorf("UNSUPPORTED_ELEMENT: transaction %q must have a startEvent", tx.ID)
	}
	return nil
}

func cancelEndSpec(ev element.EndEvent) error {
	if len(ev.CancelEventDefinitions) != 1 {
		return fmt.Errorf("not a cancel end")
	}
	other := extraCatchDefinitions(ev.EventDefinitions) - len(ev.CancelEventDefinitions)
	if other > 0 || len(ev.TimerEventDefinitions) > 0 {
		return fmt.Errorf("not a cancel end")
	}
	return nil
}

func cancelBoundarySpec(ev element.BoundaryEvent) (attachedTo string, err error) {
	if err := requireBoundaryAttach(ev); err != nil {
		return "", err
	}
	if len(ev.CancelEventDefinitions) != 1 {
		return "", fmt.Errorf("UNSUPPORTED_ELEMENT: boundaryEvent %q must be a cancel boundary", ev.ID)
	}
	other := extraCatchDefinitions(ev.EventDefinitions) - len(ev.CancelEventDefinitions)
	if other > 0 || len(ev.TimerEventDefinitions) > 0 {
		return "", fmt.Errorf("UNSUPPORTED_ELEMENT: boundaryEvent %q must be a cancel boundary", ev.ID)
	}
	if !ev.CancelActivity {
		return "", fmt.Errorf("UNSUPPORTED_ELEMENT: cancel boundary %q must be interrupting (cancelActivity=true)", ev.ID)
	}
	return strings.TrimSpace(ev.AttachedToRef), nil
}

// IsTransaction reports whether id is a BPMN transaction element.
func (d *Deployment) IsTransaction(id string) bool {
	if d == nil {
		return false
	}
	t, err := d.TypeOf(id)
	return err == nil && t == eventv1.Element_TYPE_TRANSACTION
}

// IsEmbeddedScope reports whether id is an embedded SubProcess or Transaction.
func (d *Deployment) IsEmbeddedScope(id string) bool {
	if d == nil {
		return false
	}
	t, err := d.TypeOf(id)
	if err != nil {
		return false
	}
	return (t == eventv1.Element_TYPE_SUB_PROCESS && !d.IsEventSubProcess(id)) || t == eventv1.Element_TYPE_TRANSACTION
}

// IsCancelEnd reports whether id is a cancel end event.
func (d *Deployment) IsCancelEnd(id string) bool {
	return d != nil && d.cancelEnds[id]
}

// CancelBoundaryOf returns the cancel boundary attached to transactionID.
func (d *Deployment) CancelBoundaryOf(transactionID string) (string, bool) {
	if d == nil {
		return "", false
	}
	id, ok := d.cancelBoundaries[transactionID]
	return id, ok
}

// TransactionStartEventID returns the start event of the given transaction.
func (d *Deployment) TransactionStartEventID(transactionID string) (string, error) {
	tx := d.findTransaction(&d.Process.FlowElements, transactionID)
	if tx == nil {
		for _, called := range d.calledProcesses {
			fe := called.FlowElements
			if tx = d.findTransaction(&fe, transactionID); tx != nil {
				break
			}
		}
	}
	if tx == nil {
		return "", fmt.Errorf("NOT_FOUND: transaction %q", transactionID)
	}
	if len(tx.StartEvents) == 0 {
		return "", fmt.Errorf("no startEvent in transaction %q", transactionID)
	}
	return tx.StartEvents[0].ID, nil
}

func (d *Deployment) findTransaction(fe *element.FlowElements, id string) *element.Transaction {
	for i := range fe.Transactions {
		if fe.Transactions[i].ID == id {
			return &fe.Transactions[i]
		}
		if tx := d.findTransaction(&fe.Transactions[i].FlowElements, id); tx != nil {
			return tx
		}
	}
	for i := range fe.SubProcesses {
		if tx := d.findTransaction(&fe.SubProcesses[i].FlowElements, id); tx != nil {
			return tx
		}
	}
	return nil
}

// EmbeddedScopeStartEventID returns the start event for a SubProcess or Transaction.
func (d *Deployment) EmbeddedScopeStartEventID(scopeID string) (string, error) {
	if d.IsTransaction(scopeID) {
		return d.TransactionStartEventID(scopeID)
	}
	return d.SubProcessStartEventID(scopeID)
}

func (d *Deployment) validateCancelSemantics() error {
	for endID := range d.cancelEnds {
		scope, ok := d.ScopeOf(endID)
		if !ok || !d.IsTransaction(scope) {
			return fmt.Errorf("UNSUPPORTED_ELEMENT: cancel end %q must be inside a transaction", endID)
		}
		if _, ok := d.CancelBoundaryOf(scope); !ok {
			return fmt.Errorf("UNSUPPORTED_ELEMENT: transaction %q has cancel end %q but no cancel boundary", scope, endID)
		}
	}
	for txID, bid := range d.cancelBoundaries {
		if !d.IsTransaction(txID) {
			return fmt.Errorf("UNSUPPORTED_ELEMENT: cancel boundary %q must attach to a transaction (%q)", bid, txID)
		}
	}
	return nil
}
