import React, { useCallback, useEffect, useState } from "react";
import { Accordion, Alert, Badge, Button, Card, Col, Container, ListGroup, Row, Spinner, Tab, Table, Tabs } from "react-bootstrap";
import { useNavigate, useParams } from "react-router-dom";
import "./index.css"
import SharedInfraDiagram from "../SharedInfraDiagram";
import { toEdges, toNodes } from "../SharedInfraDiagram/utils";
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

const SharedInfraView = () => {
  const navigate = useNavigate()
  const { name: sharedInfraName } = useParams()
  const [sharedInfra, setSharedInfra] = useState<any>()
  const [executions, setExecutions] = useState<any>([])
  const [selectedExecution, setSelectedExecution] = useState<boolean>()
  const [nodes, setNodes] = useState<any>([])
  const [edges, setEdges] = useState<any>([])


  const getSharedInfra = useCallback(async (name: string) => {
    const sharedInfraRes = await fetch(`http://localhost:8080/shared-infras/${name}`)
    const sharedInfra = await sharedInfraRes.json()

    setSharedInfra(sharedInfra)
  }, [])
  
  const handleReconcile = useCallback(async () => {
    const res = await fetch(`http://localhost:8080/shared-infras/${sharedInfraName}/reconcile`, {method: 'PATCH', body: JSON.stringify({})})
    const item = await res.json()
  }, [])

  useEffect(() => {
    if (!sharedInfraName) return
    
    getSharedInfra(sharedInfraName)
    interval = setInterval(() => {
      getSharedInfra(sharedInfraName)
    }, 3000)
    

    return () => clearInterval(interval)
  }, [])

  useEffect(() => {
    if (!sharedInfra) return

    if (selectedExecution) {
      setNodes(toNodes(sharedInfra?.status?.plugins, "executionNode"))
      setEdges(toEdges(sharedInfra?.status?.plugins, true))
      return
    }

    setNodes(toNodes(sharedInfra?.plugins, "default"))
    setEdges(toEdges(sharedInfra?.plugins, false))
  }, [selectedExecution, sharedInfra])

  
  return (
    <div className="shared-infra-view__content">
      <DefaultPanel
        sharedInfra={sharedInfra}
        executions={executions}
        onViewClick={() => setSelectedExecution(false)}
        onEditClick={() => navigate(`/shared-infras/${sharedInfraName}/edit`)}
        onReconcileClick={() => handleReconcile()}
        onSelectExecution={(e: any) => setSelectedExecution(true)}
      />
      {sharedInfra?.status && sharedInfra?.status?.error && (
        <Alert
          style={{position: 'fixed', top: '10px', right: '10px', left: '390px'}}
          variant="danger"
        >{sharedInfra?.status?.error}</Alert>
      )}
      <div className="shared-infra-view__diagram">
      <SharedInfraDiagram
        sharedInfra={sharedInfra}
        nodes={nodes}
        edges={edges}
      />
      </div>
     
    </div>
  )
}

const replaceBreakLines = (text: string) => text.replace(/(?:\\n|\\\\n)/g, '<br/>')

export default SharedInfraView