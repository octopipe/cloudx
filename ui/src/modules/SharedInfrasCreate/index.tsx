import React, { useCallback, useEffect, useState } from "react";
import { Accordion, Alert, Badge, Button, Card, Col, Container, ListGroup, Row, Spinner, Tab, Table, Tabs } from "react-bootstrap";
import { useParams } from "react-router-dom";
import "./index.css"
import SharedInfraDiagram from "../SharedInfraDiagram";
import { toEdges, toNodes } from "../SharedInfraDiagram/utils";

const getBadgeVariants = (status: string) => {
  if (status === "RUNNING") {
    return 'primary'
  }

  if (status === "SUCCESS") {
    return 'success'
  }


  return 'danger'
}

const SharedInfraCreate = () => {
  const { name } = useParams()
  const [sharedInfra, setSharedInfra] = useState<any>()
  const [selectedExecution, setSelectedExecution] = useState<any>()

  const getSharedInfra = useCallback(async (name: string) => {
    const res = await fetch(`http://localhost:8080/shared-infras/${name}`)
    const item = await res.json()

    setSharedInfra(item)
  }, [])


  const getExecution = useCallback(async (name: string) => {
    const res = await fetch(`http://localhost:8080/executions/${name}`)
    const item = await res.json()

    setSelectedExecution(item)
  }, [])

  useEffect(() => {
    if (!name)
      return

    const interval = setInterval(() => {
      getSharedInfra(name)
    }, 3000)

    getSharedInfra(name)
    return () => clearInterval(interval)
  }, [])
  
  return (
    <div className="shared-infra-view__content">
      <SharedInfraDiagram
        action="CREATE"
        nodes={sharedInfra?.plugins ? toNodes(sharedInfra.plugins, "defaultNode") : []}
        edges={sharedInfra?.plugins ? toEdges(sharedInfra.plugins) : []}
      />
    </div>
  )
}

const replaceBreakLines = (text: string) => text.replace(/(?:\\n|\\\\n)/g, '<br/>')

export default SharedInfraCreate