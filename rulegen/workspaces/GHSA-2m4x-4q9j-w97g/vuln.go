package main

			return true, errs
		}
		return false, errs
	case *object:
		cpy, _ := v.Map(func(k, v *Term) (*Term, *Term, error) {
			kcpy := k.Copy()
