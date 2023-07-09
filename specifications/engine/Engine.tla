---------------------- MODULE Engine ----------------------
EXTENDS Naturals

(* Conjunto de tarefas *)
VARIABLE tasks

(* Conjunto de dependências entre as tarefas *)
VARIABLE dependencies

(* Estado inicial do pipeline *)
Init == 
  /\ tasks = {1, 2, 3, 4}
  /\ dependencies = {<<1, 2>>, <<2, 3>>, <<1, 4>>}

(* Ação para executar uma tarefa *)
ExecuteTask(task) ==
  /\ task \in tasks
  /\ \A dep \in dependencies: dep[2] # task
  /\ tasks' = tasks \ {task}
  /\ dependencies' = {dep \in dependencies: dep[1] # task}
  
(* Próxima tarefa a ser executada *)
NextTask ==
  CHOOSE task \in tasks: \A dep \in dependencies: dep[2] # task

(* Comportamento do sistema *)
Next == 
  \/ ExecuteTask(NextTask)
  \/ \E task \in tasks: ExecuteTask(task)

(* Propriedade: Não há tarefas dependentes pendentes *)
NoPendingDependencies ==
  \A dep \in dependencies: dep[2] \notin tasks

(* Propriedade: O pipeline sempre termina *)
Termination ==
  tasks = {}

(* Especificação do sistema *)
Spec ==
  /\ Init
  /\ [][Next]_<<tasks, dependencies>>
  /\ NoPendingDependencies
  /\ Termination

=============================================================
