class_name enemy
extends node

# Attack range
# Yeah thats the attack range
@export range : int = 0

@export dmg : int = 0

@export speed : int = 0
@export also_speed: int = 0

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
	position.y++

func bar(delta: float) -> void:
	position.z--