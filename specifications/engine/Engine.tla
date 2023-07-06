---- MODULE Engine ----
EXTENDS TLC
CONSTANTS CPLUGINS, LPLUGINS
VARIABLES pc, applied, destroyed

vars == << pc, applied, destroyed >>

Init ==
    /\ pc = {}
    /\ applied = {}
    /\ destroyed = {}

Apply(self) ==
    \/ self \in LPLUGINS
        /\ destroyed' = destroyed \cup {self}
    \/ applied' = applied \cup {self}
    

Execute(self) ==
    /\ self \notin applied
    /\ self \notin destroyed
    /\ Apply(self)
    /\ UNCHANGED vars
    
Next == (\E self \in CPLUGINS: Execute(self))

Spec ==
    /\ Init
    /\ [][Next]_vars
    /\ WF_vars(Next)
        
\* Properties

\* All == [](\E self \in CPLUGINS)

====