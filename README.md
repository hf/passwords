Go Passwords
============

A package of properly modified password hashing algorithms that cooperate with
other Goroutines. Ideal for building web servers doing password hashing.

Please read [Go and passwords][go-and-passwords] for more information on why
you should use this package instead of:

- `golang.org/x/crypto/argon2`
- `golang.org/x/crypto/bcrypt`
- `golang.org/x/crypto/scrypt`
- `golang.org/x/crypto/pbkdf2`

## How to use

Instead of using `golang.org/x/crypto` use `github.com/hf/passwords`. You'll
find `argon2`, `bcrypt`, `scrypt`, and `pbkdf2` as supbackages.

These packages expose the same API as the native implementation, but also add
`WithContext` methods which allow you to cancel / timeout password hashing.

Inside the `metrics` package you will find some useful metrics:

- `NumOutstanding()` is the number of password hashing runs waiting in a queue.
- `DurationMovingAverage4()` is a 4-point moving average of the duration of
  password hashing runs.
- `DurationQueue()` gives you the duration likely needed to clear out the
  current queue of runs.

You can use these to implement better auto-scaling strategies, as well as give
you the ability to reject new password hashing runs if the system is too full
to handle them.

## Security

Please reach out to me directly over email sdimitrovski@gmail.com. Note that
the algorithms are as-implemented by the Go Authors, so be sure to also submit
a report there.

## License

Copyright &copy; 2009-2022 The Go Authors and Stojan Dimitrovski. All rights
reserved.

See `LICENSE` for the full text. It's a BSD-style license.

[go-and-passwords]: https://notes.stojan.me/Software/Go+and+passwords
