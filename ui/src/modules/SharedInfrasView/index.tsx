import React, { useCallback, useEffect, useState } from "react";
import { Accordion, Alert, Badge, Button, Card, Col, Container, ListGroup, Row, Spinner, Tab, Table, Tabs } from "react-bootstrap";
import { useParams } from "react-router-dom";
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

const SharedInfraView = () => {
  let interval: any
  const { name } = useParams()
  const [sharedInfra, setSharedInfra] = useState<any>()
  const [selectedExecution, setSelectedExecution] = useState<any>()
  const [nodes, setNodes] = useState<any>([])
  const [edges, setEdges] = useState<any>([])


  const getSharedInfra = useCallback(async (name: string) => {
    const res = await fetch(`http://localhost:8080/shared-infras/${name}`)
    const item = await res.json()

    setSharedInfra(item)
    setNodes(toNodes(item.plugins, "defaultNode"))
    setEdges(toEdges(item.plugins))
  }, [])


  const getExecution = useCallback(async (name: string, namespace: string) => {
    const res = await fetch(`http://localhost:8080/executions/${name}?namespace=${namespace}`)
    const item = await res.json()

    setNodes(toNodes(item?.status?.plugins || [], "executionNode"))
    setEdges(toEdges(item?.status?.plugins || []))
  }, [])

  useEffect(() => {
    if (!name)
      return
    
    if (!!selectedExecution) {
      clearInterval(interval)
      interval = setInterval(() => {
        getExecution(selectedExecution?.name, selectedExecution?.namespace)
      }, 3000)
      getExecution(selectedExecution?.name, selectedExecution?.namespace)
      return
    }

    interval = setInterval(() => {
      getSharedInfra(name)
    }, 3000)
    getSharedInfra(name)

    return () => clearInterval(interval)
  }, [selectedExecution])
  
  return (
    <div className="shared-infra-view__content">
      <DefaultPanel sharedInfra={sharedInfra} onSelectExecution={(e: any) => setSelectedExecution(e)} />
      <SharedInfraDiagram
        sharedInfra={sharedInfra}
        nodes={nodes}
        edges={edges}
      />
    </div>
  )
}

const replaceBreakLines = (text: string) => text.replace(/(?:\\n|\\\\n)/g, '<br/>')

export default SharedInfraView