// Package fate allows implementing "designing for failure" by introducing
// errors inside application logic by "tempting fate".
//
//   type Fate interface {
//      Tempt() error
//   }
//
// Fate support deterministic behavior for tests and probabilistic behavior
// for production.
//
// The canonical usage pattern:
//
//   function multiStepLogic(state state.State, fate fate.Fate, arg int) error {
//     r1, err := step1(state.Step1DB(), arg)
//     if err != nil {
//       return err
//     }
//
//     if err := fate.Temp(); err != nil {
//       return err
//     }
//
//     err := step2(state.Step2DB(), r1)
//     if err != nil {
//       return err
//     }
//
//     return fate.Temp()
//   }
//
// Increasing the tangibility and probability of errors by regularly tempting
// fate both in production and tests, re-enforces the need for idempotent
// robust failure proof code.
package fate
