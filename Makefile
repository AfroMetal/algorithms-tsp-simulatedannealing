all: julia

julia:
	# run to precompile things that will be needed
	julia simulated_annealing.jl 1 < dummy
