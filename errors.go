package errors

import (
	"log"
	"runtime/debug"
)

type TryFunc func()
type CatchFunc func(error)
type FinallyFunc func()

type TCBuilder struct {
	try     TryFunc
	catch   CatchFunc
	finally FinallyFunc
}

func Try(f TryFunc) *TCBuilder {
	return &TCBuilder{try: f}
}

func (tcb *TCBuilder) Catch(f CatchFunc) *TCBuilder {
	tcb.catch = f
	return tcb
}

func (tcb *TCBuilder) Finally(f ...FinallyFunc) {
	if len(f) > 0 {
		tcb.finally = f[0]
	}
	if tcb.finally != nil {
		defer tcb.finally()
	}
	if tcb.catch != nil {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					log.Println("stack trace", string(debug.Stack()))
					tcb.catch(err)
				} else {
					log.Panic(r)
				}
			}
		}()
	}
	tcb.try()
}

// func main() {
// 	var coba string

// 	Try(func() {
// 		fmt.Println("Trying...")
// 		// Ubah baris berikut untuk menguji error
// 		// panic(fmt.Errorf("some error occurred"))
// 		coba = "bener"
// 	}).
// 		Catch(func(e error) {
// 			log.Printf("Caught error: %v\n", e)
// 		}).
// 		Finally(func() {
// 			fmt.Println("Finally...")
// 		})

// 	fmt.Printf("coba: %s\n", coba)
// }
