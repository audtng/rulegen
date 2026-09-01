package main

	}
	cmd.Level = uint32(data[levelStart])
	timeStart := levelStart + 1
	if len(data) < timeStart+1 {
		return nil, newError("insufficient length.")
	}
	cmd.ValidMin = data[timeStart]
