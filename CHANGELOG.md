## v0.3.8 - Root Package
  * Commit: 
  * New Commands
  * ParamsMap
## v0.3.7 - Application Packages 
  * Commit: 2026-04-07 13:50
  * Fail functions
  * Conf package
  * conf.LoadEnv
  * dict.Inspect
  * Sys package and display functions
  * String JSON functions
## v0.3.6 - Schema Combos 
  * Commit: 2026-04-01 22:44
  * Schema GetOrCreate 
  * Schema GetAndLockTx 
  * Schema GetAndLockTxItems 
  * Schema GetOrCreateAndLockTx
  * Schema UpdateAndGetTx
  * MoveItemTx
## v0.3.5 - Schema Add and Edit 
  * Commit: 2026-04-01 15:16
  * Schema Insert 
  * Schema InsertRows
  * Schema UpdateFields 
  * Schema Update
## v0.3.4 - Schema Get and Group
  * Commit: 2026-04-01 11:47
  * Schema Count 
  * Schema Sum 
  * Schema Get
  * Schema GetRows
  * Update qb.Exec and ExecTx to return ds.Result
## v0.3.3 - Schema Delete and Toggle 
  * Commit: 2026-04-01 09:51
  * Schema Delete 
  * Schema Toggle 
  * Schema SetFlag
  * Schema SetFlags
## v0.3.2 - Model Package 
  * Commit: 2026-03-31 16:19
  * Model fields
  * Schema type
  * Clock format tests
  * Status code tests
## v0.3.1 - My Package 
  * Commit: 2026-03-31 15:26
  * Request type
  * Request methods
  * Clock format functions
## v0.3.0 - Rename Package
  * Commit: 2026-03-28 21:56
  * Rename package to github.com/zeroibot/pack
  * Refactor common code
## v0.2.23 - Transaction Tests
  * Commit: 2026-03-27 15:59
  * ResultChecker tests
  * Transaction tests
  * DB interfaces tests
  * DB adapters tests
  * NewSQLConnection 
## v0.2.22 - Query Execution Tests
  * Commit: 2026-03-27 11:57
  * ResultChecker functions
  * Query Execution functions
  * DeleteQuery Execute tests
  * InsertRowQuery Execute tests
  * InsertRowsQuery Execute tests
  * UpdateQuery Execute tests
## v0.2.21 - Multiple Row Query Execution Tests
  * Commit: 2026-03-26 14:57
  * DistinctValuesQuery Query tests
  * LookupQuery Lookup tests
  * TopRowQuery QueryRows tests
  * TopValueQuery QueryValues tests
  * SelectRowsQuery Query tests
  * GroupCountQuery GroupCount tests
  * GroupSumQuery GroupSum tests
  * Multiple Row Query edge cases tests
## v0.2.20 - Single Row Query Execution Tests
  * Commit: 2026-03-25 16:33
  * CountQuery Count and Exists tests
  * SumQuery Sum tests
  * ValueQuery QueryValue tests
  * SelectRowQuery QueryRow tests
  * TopRowQuery QueryRow tests
  * TopValueQuery QueryValue tests
  * Refactor prepare steps
  * Transfer tzt functions to tst
## v0.2.19 - Multiple Row Query Executions
  * Commit: 2026-03-24 09:41
  * DistinctValuesQuery Results
  * LookupQuery Results
  * TopRowQuery and TopValueQuery Results
  * SelectRowsQuery Results
  * GroupCountQuery Results
  * GroupSumQuery Results
## v0.2.18 - Single Row Query Executions
  * Commit: 2026-03-24 07:37
  * CountQuery Results
  * Add DB interfaces
  * ValueQuery Results
  * SelectRowQuery Results
  * TopRowQuery Results
  * TopValueQuery Results
  * SumQuery Results
## v0.2.17 - Aggregate Queries
  * Commit: 2026-03-23 16:19
  * Sum Query
  * Sum Query tests
  * GroupCount Query
  * GroupCount Query tests
  * GroupSum Query
  * GroupSum Query tests
## v0.2.16 - Multiple Row Query Tests
  * Commit: 2026-03-23 15:24
  * DistinctValues Query tests
  * Lookup Query tests
  * SelectRows Query tests
## v0.2.15 - Single Row Query Tests
  * Commit: 2026-03-23 14:33
  * Count Query tests
  * Value Query tests
  * SelectRow Query tests
  * TopRow Query tests
  * TopValue Query tests
## v0.2.14 - Multiple Row Queries
  * Commit: 2026-03-20 21:44
  * Refactor common query parts
  * DistinctValues Query
  * Lookup Query
  * SelectRows Query
## v0.2.13 - Single Row Queries
  * Commit: 2026-03-20 20:09
  * Count Query
  * Value Query
  * SelectRow Query
  * TopRow Query
  * TopValue Query
## v0.2.12 - Insert Query
  * Commit: 2026-03-20 13:09
  * InsertRow Query
  * InsertRows Query
  * InsertRow Query tests
  * InsertRows Query tests
## v0.2.11 - Convert QueryBuilder Tests
  * Commit: 2026-03-19 20:07
  * Convert qb tests
## v0.2.10 - Convert Data Structures Tests
  * Commit: 2026-03-19 14:11
  * Convert ds tests
## v0.2.9 - Convert List and Dict Tests
  * Commit: 2026-03-19 09:09
  * Convert list tests
  * Convert dict tests
## v0.2.8 - Tst Package
  * Commit: 2026-03-18 14:26
  * Tst package
  * Convert conv tests
  * Convert lang tests
  * Convert number tests
  * Convert str tests
  * Convert dyn tests
## v0.2.7 - Update Query
  * Commit: 2026-03-17 20:47
  * Add limitQuery
  * Embed limitQuery to DeleteQuery
  * Combine QB files
  * Update Query
  * Update Query tests
## v0.2.6 - Delete Query
  * Commit: 2026-03-17 13:39
  * Use slices.Sorted in Maps
  * Abstract Query types
  * Rename QB conditions and combos
  * Delete Query
  * Replace ds.List[any] with []any
  * Query Core tests
  * Delete Query tests
## v0.2.5 - QueryBuilder Condition Tests
  * Commit: 2026-03-17 08:54
  * QB Mock Scanner test
  * QB condition type tests
  * QB combo type tests
  * Add ds.Result type
  * Update RowReader to use ds.Result
  * QB conditions tests
  * QB combos tests
## v0.2.4 - QueryBuilder Internal Tests
  * Commit: 2026-03-16 16:31
  * Additional str and dyn tests
  * QB Rows test
  * QB internal rows and fields test
  * QB Instance tests
  * QB internal columns tests
  * QB remove column preparation
  * QB internal condition tests
## v0.2.3 - QueryBuilder Combos
  * Commit: 2026-03-15 15:04
  * QB Combo types
  * QB Combos
  * QB Combos Auto-Create Test Function
  * QB Multi Combos Test Function
## v0.2.2 - QueryBuilder Package
  * Commit: 2026-03-14 21:58
  * QB Add Type
  * QB Instance methods
  * QB Row functions
  * QB Column-Value constructors
  * QB Condition-Values builders
  * QB Condition types
  * QB Conditions
## v0.2.1 - Dyn Functions
  * Commit: 2026-03-14 09:30
  * Test dyn struct functions
  * Dyn MustGet functions
  * Dyn DerefValue and GetStructFieldTag
## v0.2.0 - Dyn Package
  * Commit: 2026-03-13 16:33
  * Dyn functions
  * Dyn struct functions
  * Test dyn type functions
  * Test dyn check functions
## v0.1.9 - List Package Tests
  * Commit: 2026-03-13 13:23
  * Test list number functions
  * Test list order functions
  * Test list check functions
  * Test list functional programming
  * Test list compare functions
  * Test list general functions
## v0.1.8 - List Method Tests
  * Commit: 2026-03-13 09:04
  * Test List type
  * Test NumList methods
  * Test List random methods
  * Test List functional methods
  * Test List check methods
  * Test List general methods
## v0.1.7 - Stack and Queue
  * Commit: 2026-03-12 20:58
  * Stack type
  * Stack methods
  * Queue type
  * Queue methods
  * Test stack methods
  * Test queue methods
## v0.1.6 - Lists
  * Commit: 2026-03-12 15:20
  * List type
  * List methods
  * Equality functions
  * List check functions
  * List order functions
  * List number functions
  * List compare queries and functions
  * List functional programming
  * Map-to-struct functions
## v0.1.5 - Sets
  * Commit: 2026-03-11 16:30
  * Set type
  * Set methods
  * Map tests
  * Set tests
  * Map update functions
  * Map zips and transforms
  * Map group functions
## v0.1.4 - Maps
  * Commit: 2026-03-10 16:03
  * Map type
  * Map functions
  * Custom map types
  * Map sorted functions
## v0.1.3 - Option, Tuples, Ranges
  * Commit: 2026-03-10 10:45
  * Option type
  * Tuple types: Tuple2, Tuple3, Tuple4
  * Homogeneous Tuple types: Pair, Triple, Quad
  * Range type: Normal, Reversed
## v0.1.2 - Strings and Lang Package
  * Commit: 2026-03-09 16:55
  * Lang package
  * String functions
  * String parts functions
  * String builder
## v0.1.1 - Numbers and Conversions
  * Commit: 2026-03-09 14:58
  * Number package
  * Conv package: bool, float, int
  * Number parse functions
  * Number format functions
  * Number math functions
## v0.1.0 - Data Types
  * Commit: 2026-03-01 18:40
  * ds.List, Map, Set
  * ds.Range
  * List methods
  * Map methods
  * Set methods
  * Conv package
  * Transfer methods to list package
  * Transfer methods to dict package
  * String functions