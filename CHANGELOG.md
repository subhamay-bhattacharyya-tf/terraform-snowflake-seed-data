# Changelog

All notable changes to this project will be documented in this file.

## 1.0.0 (2026-02-25)

### ⚠ BREAKING CHANGES

* Complete module rewrite from Snowflake warehouse management to seed data functionality

- Replace warehouse_configs with seed configuration object
- Add environment-based safety controls (blocked_environments)
- Support inline SQL (sql_text) or external script files (script_path)
- Add re-run controls via seed_version and rerun_on_script_change
- Update all examples, tests, and documentation

### Features

* convert warehouse module to seed data module ([a45b290](https://github.com/subhamay-bhattacharyya-tf/terraform-snowflake-seed-data/commit/a45b29024fa0bdee4253aa62a52e2ba880e8d449))
* **snowflake:** upgrade provider to 0.89.0 and migrate to snowflake_execute ([0722f6e](https://github.com/subhamay-bhattacharyya-tf/terraform-snowflake-seed-data/commit/0722f6e972aa6828ca2cb0a92bfefe02950ada6c))

## [unreleased]

### 🚀 Features

- [**breaking**] Convert warehouse module to seed data module
- *(snowflake)* Upgrade provider to 0.89.0 and migrate to snowflake_execute

### 📚 Documentation

- Update CHANGELOG.md [skip ci]

### 🧪 Testing

- *(seed)* Convert SQL strings to single-line format
