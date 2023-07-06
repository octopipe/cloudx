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
  const [selectedExecution, setSelectedExecution] = useState<any>()
  const [currentExecution, setCurrentExecution] = useState<any>()
  const [nodes, setNodes] = useState<any>([])
  const [edges, setEdges] = useState<any>([])


  const getSharedInfra = useCallback(async (name: string) => {
    const sharedInfraRes = await fetch(`http://localhost:8080/shared-infras/${name}`)
    const sharedInfra = await sharedInfraRes.json()


    const executionsRes = await fetch(`http://localhost:8080/shared-infras/${name}/executions`)
    const executions = await executionsRes.json()

    setSharedInfra(sharedInfra)
    setExecutions(executions?.items)
    setNodes(toNodes(sharedInfra.plugins, "default"))
    setEdges(toEdges(sharedInfra.plugins, false))
  }, [])


  const getExecution = useCallback(async (name: string, namespace: string) => {
    const res = await fetch(`http://localhost:8080/shared-infras/${sharedInfraName}/executions/${name}?namespace=${namespace}`)
    const item = await res.json()
    setCurrentExecution(item)
    setNodes(toNodes(item?.status?.plugins || [], "executionNode"))
    setEdges(toEdges(item?.status?.plugins || [], true))
  }, [])

  useEffect(() => {
    if (!sharedInfraName)
      return
    
    if (!!selectedExecution) {
      clearInterval(interval)
      interval = setInterval(() => {
        getExecution(selectedExecution?.name, selectedExecution?.namespace)
      }, 3000)
      getExecution(selectedExecution?.name, selectedExecution?.namespace)
      return
    }

    clearInterval(interval)
    interval = setInterval(() => {
      getSharedInfra(sharedInfraName)
    }, 3000)
    getSharedInfra(sharedInfraName)

    return () => clearInterval(interval)
  }, [selectedExecution])
  
  return (
    <div className="shared-infra-view__content">
      <DefaultPanel
        sharedInfra={sharedInfra}
        executions={executions}
        onViewClick={() => setSelectedExecution(null)}
        onEditClick={() => navigate(`/shared-infras/${sharedInfraName}/edit`)}
        onSelectExecution={(e: any) => setSelectedExecution(e)}
      />
      {currentExecution && currentExecution?.status?.error && (
        <Alert
          style={{position: 'fixed', top: '10px', right: '10px', left: '390px'}}
          variant="danger"
        >{currentExecution?.status?.error}</Alert>
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