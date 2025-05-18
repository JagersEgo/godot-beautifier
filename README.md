# About Godot Beautifier

**Godot Beautifier** is a tool designed to automatically format GDScript code to improve readability and maintain consistency across Godot Engine projects.
**Be careful, this project is new and may break code if it hits an edge case I havent thought of, this will back up your code but please back it up yourself too**

## Purpose

GDScript formatting can be inconsistent, especially when multiple developers contribute or when code is quickly written. This beautifier:

- Enforces GDScript style in these areas:
  - Code block order (https://docs.godotengine.org/en/stable/tutorials/scripting/gdscript/gdscript_styleguide.html#code-order)
  - 2x new lines between function declarations
  - No trailing newlines/indentation

## Features
- Supports the Godot 4.x GDScript syntax
- Preserve variable blocks and comments tied to code blocks
- Puts a backup before making changes in your OS's temp directory
- Command-line interface
- Processes single files or entire godot project

## Installation/Usage

Clone the repository and run the tool with go
`go run ./ {PATH TO PROJECT}`

## Example
Before (Bad layout and spacing):
```Python
extends node

# Move
func _process(delta: float) -> void:
	position.x += 1
	var foo = "bar"
	print(banana)

# Fruits block
@onready var fruit0 = "apple"
@onready var fruit1 = "banana"
@onready var fruit2 = "cherry"
@onready var fruit3 = "date"

# Attack range
# Yeah thats the attack range
@export var range : int = 0

@export var dmg : int = 0

@export var speed : int = 0
@export var also_speed: int = 0


class_name enemy

func foo(delta: float) -> void:
	position.y += 1

func bar(delta: float) -> void:
	position.z -= 1
```

After:
```Python
class_name enemy
extends node

# Attack range
# Yeah thats the attack range
@export range : int = 0

@export dmg : int = 0

@export var speed : int = 0
@export var also_speed: int = 0

# Fruits block
@onready var fruit0 = "apple"
@onready var fruit1 = "banana"
@onready var fruit2 = "cherry"
@onready var fruit3 = "date"

# Move
func _process(delta: float) -> void:
	position.x += 1
	var foo = "bar"
	print(banana)


func foo(delta: float) -> void:
	position.y += 1


func bar(delta: float) -> void:
	position.z -= 1
```
## To do
- Full adherence to style guidelines
- Add GUI version
- Able to change tscn internal scripts
