---------------------- MODULE Engine ----------------------
EXTENDS TLC, Naturals, Sequences, FiniteSets

CONSTANTS tasks, dependencies, last_execution_tasks, last_execution_dependencies
VARIABLES applied, destroyed

vars == << applied, destroyed >>

Init ==
  /\ applied = {}
  /\ destroyed = {}

ApplyTask(task) ==
  /\ task \notin applied
  /\ (dependencies[task] = {}) \/ (\A dep \in dependencies[task]: dep \in applied)
  /\ applied' = applied \cup {task}
  /\ UNCHANGED destroyed

DestroyTask(last_task) ==
  /\ last_task \notin tasks
  /\ last_task \notin destroyed
  /\ destroyed' = destroyed \cup {last_task}
  /\ UNCHANGED applied

Destroy == \E self \in last_execution_tasks: DestroyTask(self)
Apply == \E self \in tasks: ApplyTask(self)

Next ==
  \/ Destroy
  \/ Apply
  \/ UNCHANGED vars

Spec ==
  /\ Init
  /\ [][Next]_vars
  /\ WF_vars(Next)


\* Properties
Agreement == <>(\A self \in tasks: self \in applied)

=============================================================
