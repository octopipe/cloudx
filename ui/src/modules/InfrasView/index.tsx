import React, { useCallback, useEffect, useState } from "react";
import { Accordion, Alert, Badge, Button, Card, Col, Container, ListGroup, Row, Spinner, Tab, Table, Tabs } from "react-bootstrap";
import { useNavigate, useParams } from "react-router-dom";
import "./index.css"
import InfraDiagram from "../InfraDiagram";
import { toEdges, toNodes } from "../InfraDiagram/utils";
import DefaultPanel from "./DefaultPanel";

const getBadgeVariants = (status: string) => {
  if (status === "RUNNING") {
    return 'primary'
  }

  if (status === "SUCCESS") {
    return 'success'
  }


  return 'danger'
}

let interval: any

const InfraView = () => {
  const navigate = useNavigate()
  const { name: infraName } = useParams()
  const [infra, setInfra] = useState<any>()
  const [executions, setExecutions] = useState<any>([])
  const [selectedExecution, setSelectedExecution] = useState<boolean>()
  const [nodes, setNodes] = useState<any>([])
  const [edges, setEdges] = useState<any>([])


  const getInfra = useCallback(async (name: string) => {
    const infraRes = await fetch(`http://localhost:8080/infra/${name}`)
    const infra = await infraRes.json()

    setInfra(infra)
  }, [])
  
  const handleReconcile = useCallback(async () => {
    const res = await fetch(`http://localhost:8080/infra/${infraName}/reconcile`, {method: 'PATCH', body: JSON.stringify({})})
    const item = await res.json()
  }, [])

  useEffect(() => {
    if (!infraName) return
    
    getInfra(infraName)
    interval = setInterval(() => {
      getInfra(infraName)
    }, 3000)
    

    return () => clearInterval(interval)
  }, [])

  useEffect(() => {
    if (!infra) return

    if (selectedExecution) {
      setNodes(toNodes(infra?.status?.tasks, "executionNode"))
      setEdges(toEdges(infra?.status?.tasks, true))
      return
    }

    setNodes(toNodes(infra?.tasks, "default"))
    setEdges(toEdges(infra?.tasks, false))
  }, [selectedExecution, infra])

  
  return (
    <div className="shared-infra-view__content">
      <DefaultPanel
        infra={infra}
        executions={executions}
        onViewClick={() => setSelectedExecution(false)}
        onEditClick={() => navigate(`/infra/${infraName}/edit`)}
        onReconcileClick={() => handleReconcile()}
        onSelectExecution={(e: any) => setSelectedExecution(true)}
      />
      {infra?.status && infra?.status?.error && (
        <Alert
          style={{position: 'fixed', top: '10px', right: '10px', left: '390px'}}
          variant="danger"
        >{infra?.status?.error}</Alert>
      )}
      <div className="shared-infra-view__diagram">
      <InfraDiagram
        infra={infra}
        nodes={nodes}
        edges={edges}
      />
      </div>
     
    </div>
  )
}

const replaceBreakLines = (text: string) => text.replace(/(?:\\n|\\\\n)/g, '<br/>')

export default InfraView