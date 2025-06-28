# DBaker

baking...

- [x] add intermediate representation step
  - [x] write representation during introspect
  - [x] read representation during generate
  - [ ] define & validate json schema in both cases
- [x] add command line interface using cobra
- [ ] improve error handle / error message reporting
- [ ] add support for insert batching
- [ ] add ability to anotate table fields
- [ ] add support for foreign keys (1:1, 1:N, N:M)
- [ ] add support for composite primary keys
- [ ] add support for composite foregin keys
- [ ] paralelise value generating (v1 - just use number of CPUs and split the workload)
- [ ] paralelise value generating (v2 - add configurable number of parallel generators and db connections - batching X goroutines)
- [ ] add test suite (integration test with live postgres via docker & test containers)
- [ ] add charmbracelet to improve the user experience while waiting
