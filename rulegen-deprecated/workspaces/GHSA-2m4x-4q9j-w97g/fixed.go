package main

			return true, errs
		}
		return false, errs
	case Call:
		ref := v[0]
		WalkVars(ref, func(v Var) bool {
			if gv, ok := stack.Declared(v); ok && !gv.Equal(v) {
				// We will rewrite the ref of a function call, which is never ok since we don't have first-class functions.
				errs = append(errs, NewError(CompileErr, term.Location, "called function %s shadowed", ref))
				return true
			}
			return false
		})
		return false, errs
	case *object:
		cpy, _ := v.Map(func(k, v *Term) (*Term, *Term, error) {
			kcpy := k.Copy()
