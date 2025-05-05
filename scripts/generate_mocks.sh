#!/bin/bash

go install go.uber.org/mock/mockgen@v0.5.1

rm -r ./internal/competition/mocks/*

mockgen -destination=internal/competition/mocks/observer.go -package=mock_observer github.com/AleksandrMatsko/yadro-biathlon/internal/competition Observer
