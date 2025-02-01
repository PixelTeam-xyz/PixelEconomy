package main

import (
    "msg"
    "os"
)

func Except(args ...any) {
    var err error
    var exitCode *int
    var f string

    for _, arg := range args {
        switch v := arg.(type) {
        case error:
            err = v
        case int:
            exitCode = &v
        case string:
            f = v
        }
    }

    if f == "" {
        f = "%s"
    }

    if err != nil {
        if exitCode != nil {
            msg.Errorf(f, err.Error())
            os.Exit(*exitCode)
        } else {
            msg.Fatalf(f, err.Error())
        }
    }
}
