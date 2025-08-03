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
- [ ] add support for additional/missing PostgreSQL types (e.g., serial, bigserial, numeric, money, json, jsonb, bytea, inet, cidr, macaddr, bit, bit varying, interval, arrays, enums, geometric, range, xml, OID types)

## Example: Running DBaker against test tables

1. Start a local Postgres instance and initialize it with `test/init/init-tables.sql` (already done by docker-compose).
2. Run the introspection step to generate the recipe (replace credentials as needed):

```sh
dbaker introspect \
  --host localhost \
  --port 5432 \
  --database postgres \
  --username postgres \
  --password password \
  --tables public.users \
  --tables public.groups \
  --tables public.numbers_test \
  --tables public.special_test \
  --tables public.datetime_test
```

3. Generate and insert fake data into the tables:

```sh
dbaker generate \
  --host localhost \
  --port 5432 \
  --database postgres \
  --username postgres \
  --password password \
  --size 100
```

This will introspect the schema and populate the supported test tables with fake data.
