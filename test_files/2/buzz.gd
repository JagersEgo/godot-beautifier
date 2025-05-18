extends RigidBody2D

@export var speed: float = 100.0
@export var turn_speed_curve: Curve

@export var rules: Array[MovementRule]

@onready var avoidance_shapecast: ShapeCast2D = $avoidance_shapecast

func _ready() -> void:
	linear_velocity = Vector2(randf_range(-1, 1), randf_range(-1, 1)).normalized() * speed
	setup_avoidance_shapecast()

func setup_avoidance_shapecast() -> void:
	avoidance_shapecast.collide_with_areas = true
	avoidance_shapecast.collide_with_bodies = false

func _physics_process(delta: float) -> void:
	calculate_velocity()

func calculate_velocity() -> void:
	var result_vector: Vector2
	var magnitude: float
	for i in rules.size():
		if !rules[i].enabled:
			continue
		result_vector += calculate_force(i)
		magnitude = max(magnitude, result_vector.length())

	result_vector = result_vector.normalized()
	result_vector *= magnitude

	if result_vector == Vector2.ZERO:
		result_vector = linear_velocity.normalized()

	var turn_speed: float = turn_speed_curve.sample(result_vector.length())

	var desired_rotation = result_vector.angle()

	rotation = lerp_angle(rotation, desired_rotation, turn_speed)
	linear_velocity = Vector2.RIGHT.rotated(rotation) * speed
func get_boids_in_range(radius: float) -> Array[RigidBody2D]:
	var location: Vector2 = global_position
	var objs: Array[RigidBody2D]
	var space_state = get_world_2d().direct_space_state

	# Create a circular shape
	var circle_shape = CircleShape2D.new()
	circle_shape.radius = radius

	# Set up query parameters
	var shape_query = PhysicsShapeQueryParameters2D.new()
	shape_query.shape = circle_shape
	shape_query.transform = Transform2D(0, location)  # Position the shape at location
	shape_query.collide_with_areas = false
	shape_query.collide_with_bodies = true

	# Perform the query
	var results = space_state.intersect_shape(shape_query)

	# check result is in group and add it to the list
	for obj in results:
		if obj["collider"] == self:
			continue
		if obj["collider"].is_in_group("Boid2"):
			objs.append(obj["collider"])

	return objs
func get_closest_boid_in_range(range: float) -> RigidBody2D:
	var closest = null
	var min_dist = 2 * range
	for obj in get_boids_in_range(range):
		if global_position.distance_to(obj.global_position) < min_dist:
			min_dist = global_position.distance_to(obj.global_position)
			closest = obj

	return closest
func calculate_force(index: int) -> Vector2:
	var rule: MovementRule = rules[index]
	var boids: Array[RigidBody2D]