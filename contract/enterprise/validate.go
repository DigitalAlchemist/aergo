package enterprise

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/aergoio/aergo/state"
	"github.com/aergoio/aergo/types"
)

const SetConf = "setConf"
const AppendConf = "appendConf"
const RemoveConf = "removeConf"
const AppendAdmin = "appendAdmin"
const RemoveAdmin = "removeAdmin"
const EnableConf = "enableConf"
const DisableConf = "disableConf"

var ErrTxEnterpriseAdminIsNotSet = errors.New("admin is not set")

func ValidateEnterpriseTx(tx *types.TxBody, sender *state.V,
	scs *state.ContractState) (*EnterpriseContext, error) {
	var ci types.CallInfo
	if err := json.Unmarshal(tx.Payload, &ci); err != nil {
		return nil, err
	}
	context := &EnterpriseContext{Call: &ci}
	switch ci.Name {
	case AppendAdmin, RemoveAdmin:
		if len(ci.Args) != 1 { //args[0] : key, args[1:] : values
			return nil, fmt.Errorf("invalid arguments in payload for SetAdmin: %s", ci.Args)
		}

		arg := ci.Args[0].(string)
		context.Args = append(context.Args, arg)

		address := types.ToAddress(arg)
		if len(address) == 0 {
			return nil, fmt.Errorf("invalid arguments[0]: %s", ci.Args[0])
		}

		admins, err := checkAdmin(scs, sender.ID())
		if err != nil &&
			err != ErrTxEnterpriseAdminIsNotSet {
			return nil, err
		}

		context.Admins = admins
		if ci.Name == AppendAdmin && context.IsAdminExist(address) {
			return nil, fmt.Errorf("already exist admin: %s", ci.Args[0])
		} else if ci.Name == RemoveAdmin && !context.IsAdminExist(address) {
			return nil, fmt.Errorf("admins is not exist : %s", ci.Args[0])
		}

	case SetConf, AppendConf, RemoveConf:
		if len(ci.Args) <= 1 { //args[0] : key, args[1:] : values
			return nil, fmt.Errorf("invalid arguments in payload for setConf: %s", ci.Args)
		}
		if ci.Args[0] == "admin" {
			return nil, fmt.Errorf("not allowed key : %s", ci.Args[0])
		}
		for _, v := range ci.Args {
			arg, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("not string in payload for setConf : %s", ci.Args)
			}
			if strings.Contains(arg, "\\") {
				return nil, fmt.Errorf("not allowed charactor in %s", arg)
			}
			context.Args = append(context.Args, arg)
		}
		admins, err := checkAdmin(scs, sender.ID())
		if err != nil {
			return nil, err
		}
		context.Admins = admins
	case EnableConf:
		if len(ci.Args) != 2 { //args[0] : key, args[1] : true/false
			return nil, fmt.Errorf("invalid arguments in payload for enableConf: %s", ci.Args)
		}
		arg0, ok := ci.Args[0].(string)
		if !ok {
			return nil, fmt.Errorf("not string in payload for enableConf : %s", ci.Args)
		}
		if ci.Args[0] == "admin" {
			return nil, fmt.Errorf("not allowed key : %s", ci.Args[0])
		}
		context.Args = append(context.Args, arg0)
		_, ok = ci.Args[1].(bool)
		if !ok {
			return nil, fmt.Errorf("not bool in payload for enableConf : %s", ci.Args)
		}
		admins, err := checkAdmin(scs, sender.ID())
		if err != nil {
			return nil, err
		}
		context.Admins = admins
	default:
		return nil, fmt.Errorf("unsupported call %s", ci.Name)
	}
	return context, nil
}

func checkAdmin(scs *state.ContractState, address []byte) ([][]byte, error) {
	admins, err := getAdmins(scs)
	if err != nil {
		return nil, fmt.Errorf("could not get admin in enterprise contract")
	}
	if admins == nil {
		return nil, ErrTxEnterpriseAdminIsNotSet
	}
	if i := bytes.Index(bytes.Join(admins, []byte("")), address); i == -1 && i%types.AddressLength != 0 {
		return nil, fmt.Errorf("admin address not matched")
	}
	return admins, nil
}
