package common

import "log"

func ExitWhenError(e error) {
    if e != nil {
        log.Fatal(e)
    }
}
