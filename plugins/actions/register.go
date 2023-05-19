package actions

// Simple helper function to register an action
func RegisterActions(actions ...*Action) {
	Actions = append(Actions, actions...)
}

type actionList []*Action

func (l actionList) Find(name string) (*Action, bool) {
	for _, action := range l {
		if action.Name == name {
			return action, true
		}
	}

	return nil, false
}

var Actions actionList
