---- MODULE MCEngine ----
EXTENDS TLC, Engine

MCTasks == {1, 2, 3, 4}
MCDependencies == <<{}, {}, {1, 2}, {3}>>
MCLastExecutionTasks == {1, 2, 3, 4, 5, 6}
MCLastExecutionDependencies == <<{}, {}, {1, 2}, {3}, {2}, {1, 5}>>
====