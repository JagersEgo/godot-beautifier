extends RigidBody2D

@onready var avoidance_shapecast: ShapeCast2D = $avoidance_shapecast

@export var speed: float = 100.0
@export var turn_speed_curve: Curve

@export var rules: Array[MovementRule]

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

	var final_direction: Vector2 = Vector2.ZERO
	var final_position: Vector2 = Vector2.ZERO
	var final_vector = Vector2.ZERO

	if rule.selection_type == MovementRule.SelectionType.AllInArea:
		boids = get_boids_in_range(rule.activation_radius)
		if boids.size() > 0:
			for boid in boids:
				final_direction += Vector2.RIGHT.rotated(boid.rotation)
				final_position += boid.global_position
			final_position /= float(boids.size())

	elif rule.selection_type == MovementRule.SelectionType.Closest:
		boids.append(get_closest_boid_in_range(rule.activation_radius))
		if boids.size() > 0 and boids[0] != null:
			for boid in boids:
				final_direction += Vector2.RIGHT.rotated(boid.rotation)
				final_position += boid.global_position
			final_position /= float(boids.size())

	elif  rule.selection_type == MovementRule.SelectionType.Raycast:
		avoidance_shapecast.target_position = Vector2(rule.activation_radius, 0)
		if avoidance_shapecast.is_colliding():
			final_direction = avoidance_shapecast.get_collision_normal(0)
			final_position = avoidance_shapecast.get_collision_point(0)

	elif rule.selection_type == MovementRule.SelectionType.Random:
		final_direction = Vector2(randf_range(-1, 1), randf_range(-1, 1)).normalized()
		final_position = final_direction

	var normalised_distance = 1 - global_position.distance_to(final_position) / rule.activation_radius

	final_position -= global_position

	if rule.force_direction == MovementRule.ForceType.AverageDirection:
		final_vector = final_direction.normalized()
	elif rule.force_direction == MovementRule.ForceType.AveragePosition:
		final_vector = final_position.normalized()

	if rule.invert:
		final_vector = -final_vector

	final_vector = final_vector.normalized() * rule.force_strength.sample(normalised_distance)

	return final_vector
