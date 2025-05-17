# About Godot Beautifier

**Godot Beautifier** is a tool designed to automatically format GDScript code to improve readability and maintain consistency across Godot Engine projects.
**Be careful, this project is new and may break code if it hits an edge case I havent thought of, this will back up your code but please back it up yourself too**

## Purpose

GDScript formatting can be inconsistent, especially when multiple developers contribute or when code is quickly written. This beautifier:

- Enforces GDScript style in these areas:
  - code block order (https://docs.godotengine.org/en/stable/tutorials/scripting/gdscript/gdscript_styleguide.html#code-order)
  - 2x new lines between function declarations
  - No trailing newlines/indentation

## Features

- Supports the Godot 4.x GDScript syntax
- Puts a backup before making changes in your OS's temp directory
- Command-line interface
- Processes single files or entire godot project

## Installation/Usage

Clone the repository and run the tool with go
`go run ./ {PATH TO PROJECT}`

## To do
- Full adherence to style guidelines
- Add GUI version
- Able to change tscn internal scripts
