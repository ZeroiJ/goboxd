# Plug-and-Play Languages

goboxd uses a purely configuration-driven approach to language support. All supported languages are defined in `languages.yaml`.

## Adding a New Language
You can add a new language with **zero Go code changes**.
1. Open `languages.yaml`.
2. Add a new block for the language specifying `id`, `name`, `source_filename`, `run.cmd`, and `run.args`.
3. If the language is compiled, provide a `build.cmd` and `build.args` block, along with a `flag_allowlist`.
4. Ensure the relevant compiler/runtime is installed in the runtime stage of the `Dockerfile`.

## Template Variables
The `args` arrays support the following template expansions:
- `{{source}}`: Path to the written source file.
- `{{artifact}}`: Path to the compiled artifact.
- `{{flags}}`: Dynamically injects allowed build/run flags provided in the HTTP request.

## Security
Flags passed by the user must match the `flag_allowlist` (which supports exact matches or `*` wildcard suffixes like `-std=*`). Any unallowed flags will result in an HTTP 400 error.
