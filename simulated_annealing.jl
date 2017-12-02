# using Distances

const start = time()

const DEBUG = length(ARGS) >=1 && ARGS[1] == "1"

function distance(solution::Array{Int32})
    cost = 0
    for i in 1:length(solution)-1
        cost += d[solution[i],solution[i+1]]
    end
    return cost
end

tic()
tic()
# file lines iterator object
# input = eachline(STDIN)

# state = start(input)
# read data size N
const n = parse(Int32, readline(STDIN))
const k = max(0.01*n, 25.0)

# 2xN node's coordinates matrix
data = Array{Tuple{Float64, Float64}}(n)

# read file data to Data matrix
for i in 1:n
    v = readline(STDIN)
    (i_v, x_v, y_v) = split(v)
    data[parse(Int32, i_v)] = (parse(Float64, x_v), parse(Float64, y_v))
end
# calculate nodes distances as NxN matrix

maxtime = parse(Float64, readline(STDIN))

# d = Array{Float64}(n, n)
# pairwise!(d, Euclidean(1e-12), data)

d = [hypot(x2-x1, y2-y1) for (x2, y2) in data, (x1, y1) in data]

currentsol = Array{Int32}(n+1)
currentsol[1] = currentsol[n+1] = Int32(1)

for i in 1:n
    d[i,i] = Inf
end

d1 = copy(d)
j = 1
for i in 2:length(currentsol)-1
    currentsol[i] = findfirst(d1[j,:].==min(d1[j,:]...))
    j = currentsol[i]
    d1[j,1] = Inf
    for visited_node in 1:i
        d1[j, currentsol[visited_node]] = Inf
    end
end

currentcost = distance(currentsol)

optimumsol = copy(currentsol)
optimumcost = currentcost

####### SETTINGS #######
# how many times to iterate (minimum)
const minItrs     = 10000000

# fraction of melting point for starting temperature
const meltPointF  = 0.7

# fraction of melting point for ending temperature
const targetTempF = 0.01

# fraction of stagnant minItrs allowed before reheating
const stagItrsF   = 0.1
########################

function iterate(t::Float64, xsol::Array{Int32}, xcost::Float64)
	ysol = copy(xsol)
    ycost = xcost

	t1 = rand(2:n-2)
	t2 = rand(t1+2:n)

    ycost = ycost -
    		d[ysol[t1-1], ysol[t1]] -
    		d[ysol[t1], ysol[t1+1]] +
    		d[ysol[t1-1], ysol[t2]] +
    		d[ysol[t2], ysol[t1+1]] -
    		d[ysol[t2-1], ysol[t2]] -
    		d[ysol[t2], ysol[t2+1]] +
    		d[ysol[t2-1], ysol[t1]] +
    		d[ysol[t1], ysol[t2+1]]

	ysol[t1], ysol[t2] = ysol[t2], ysol[t1]

	if ycost <= xcost
		return ysol, ycost
	elseif rand() < exp((xcost-ycost)/t)
		return ysol, ycost
	end
	return xsol, xcost
end

function iterateP(p::Float64, xsol::Array{Int32}, xcost::Float64)
	ysol = copy(xsol)
    ycost = xcost

	t1 = rand(2:n-2)
    t2Max = round(min(t1+k, n))
	t2 = rand(Int64(t1+2):Int64(t2Max))

    ycost = ycost -
    		d[ysol[t1-1], ysol[t1]] -
    		d[ysol[t1], ysol[t1+1]] +
    		d[ysol[t1-1], ysol[t2]] +
    		d[ysol[t2], ysol[t1+1]] -
    		d[ysol[t2-1], ysol[t2]] -
    		d[ysol[t2], ysol[t2+1]] +
    		d[ysol[t2-1], ysol[t1]] +
    		d[ysol[t1], ysol[t2+1]]

	ysol[t1], ysol[t2] = ysol[t2], ysol[t1]

	if ycost <= xcost
		return ysol, ycost
	elseif rand() < p
		return ysol, ycost
	end
	return xsol, xcost
end

const initialcost = currentcost
if DEBUG
    println(string("initial cost: ", initialcost))
    println(string("min iterations: ", minItrs))
    println(string("initial time: ", toq(), " seconds"))
end

testsol = copy(currentsol)
testcost = currentcost
minT = 10.0e10
maxT = 0

for i in 1:max(0.01*minItrs, 2.0)
    testsol, testcost = iterateP(0.001, testsol, testcost)

    minT = min(minT, testcost)
    maxT = max(maxT, testcost)
end

const meltPoint = (maxT - minT) * 10.0^(-1.0*log10(n))

const t0 = meltPoint * meltPointF

if DEBUG
    @printf("T0 = %.2f\n", t0)
end
t = t0

const tDecay = targetTempF ^ (1.0/minItrs)

itr = 1
optItr = 1
stagItrs = 0

function printresult()
    if DEBUG
        @printf("T = %.2f\n", t)
    	@printf("Optimum: %.5f\n", optimumcost)
    	@printf("%.2f%% improvement\n", 100.0-100*optimumcost/initialcost)
    	@printf("after %.2f%% of iterations\n", 100*optItr/itr)
    	@printf("done %.2f%% of minimum iterations\n", 100*itr/minItrs)
    	toc()
    else
        println(STDOUT, optimumcost)
        for city in optimumsol
            print(STDERR, city)
            print(STDERR, " ")
        end
        println(STDERR)
    end
end

atexit(printresult)

while t > t0*targetTempF

    if time() - start > maxtime - 1.0
        exit()
    end

	currentsol, currentcost = iterate(t, currentsol, currentcost)
	if currentcost < optimumcost
		optimumsol = copy(currentsol)
        optimumcost = currentcost
		optItr = itr
	else
		stagItrs += 1
	end

	if stagItrs == Int(stagItrsF*minItrs) && itr <= minItrs
		stagItrs = 0
        if DEBUG
            println("reheating")
        end
		t *= 1.0 + 0.7*(minItrs-itr)/minItrs
	end

	t *= tDecay

	if DEBUG && itr%(minItrs/5) == 0
		@printf("T = %.2f\n", t)
		@printf("cost: %.5f\n", currentcost)
		@printf("optimum: %.5f\n", optimumcost)
	end
    itr += 1
end
