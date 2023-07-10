import React, { memo, useCallback, useEffect, useState } from "react";
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

const InfraEditor = memo(() => {
  const navigate = useNavigate()
  const { name } = useParams()
  const [infra, setInfra] = useState()
  const [nodes, setNodes] = useState([]);
  const [edges, setEdges] = useState([]);
  const [tasks, setTasks] = useState<any>([])

  const getInfra = useCallback(async (name: string) => {
    const infraRes = await fetch(`http://localhost:8080/infra/${name}`)
    const infra = await infraRes.json()

    setInfra(infra)
    setNodes(toNodes(infra.tasks, "default"))
    setEdges(toEdges(infra.tasks, false))
  }, [])

  useEffect(() => {
    if (!name)
      return

    getInfra(name)
  }, [name])

  const createInfra = useCallback(async (infra: any) => {
    const res = await fetch(`http://localhost:8080/infra`, {
      method: 'POST',
      body: JSON.stringify({
        ...infra,
        tasks,
      })
    })
    const created = await res.json()
    navigate(`/infra/${infra?.name}`)
  }, [tasks])

  const handleDiagramChanges = useCallback((nodes: any, edges: any) => {
    let dict: any = {}
    for(let i = 0; i < nodes.length; i++) {
      dict[nodes[i].id] = nodes[i]?.data?.name
    }

    const newTasks = nodes?.map((node: any) => {
      return {
        name: node?.data?.name,
        ref: node?.data?.ref,
        type: node?.data?.type,
        depends: edges.filter((e: any) => e.target === node.id).map((e: any) => dict[e.source]),
        inputs: node?.data?.inputs,
        outputs: [],
      }
    })

    setTasks(newTasks)
  }, [setTasks])

  useEffect(() => {
    handleDiagramChanges(nodes, edges)
  }, [nodes, edges])

  useEffect(() => {
    console.log(tasks)
  }, [tasks])

  
  return (
    <div className="shared-infra-create__content">
      <DefaultPanel
        infra={infra}
        onSave={createInfra}
        goToView={() => navigate(`/infra/${name}`)}
      />
      <div className="shared-infra-view__diagram">
      <InfraDiagram
        action="CREATE"
        nodes={nodes}
        edges={edges}
        onChangeDiagram={handleDiagramChanges}
      />
      </div>
    </div>
  )
})

export default InfraEditor