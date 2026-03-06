#!/bin/bash

set -e

mkdir -p dataset/easy
mkdir -p dataset/medium
mkdir -p dataset/hard

seed=1

generate_set () {

    difficulty=$1
    size=$2

    echo "Generating $difficulty puzzles"

    for i in $(seq -w 1 100)
    do
        meta="dataset/$difficulty/meta_$i.txt"

        # generate metadata
        ./tools/generate $size $size $seed > "$meta"

        # render images
        go run tools/render.go < "$meta"

        # move generated images
        mv nonogram_problem.png dataset/$difficulty/puzzle_$i.png
        mv nonogram_solution.png dataset/$difficulty/solution_$i.png

        seed=$((seed+1))
    done
}

generate_set easy 7
generate_set medium 12
generate_set hard 20

echo "Dataset generation complete"