import React, { memo, useCallback, useEffect, useState } from "react";
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

const SharedInfraEditor = memo(() => {
  const navigate = useNavigate()
  const { name } = useParams()
  const [sharedInfra, setSharedInfra] = useState()
  const [nodes, setNodes] = useState([]);
  const [edges, setEdges] = useState([]);
  const [plugins, setPlugins] = useState<any>([])

  const getSharedInfra = useCallback(async (name: string) => {
    const sharedInfraRes = await fetch(`http://localhost:8080/shared-infras/${name}`)
    const sharedInfra = await sharedInfraRes.json()

    setSharedInfra(sharedInfra)
    setNodes(toNodes(sharedInfra.plugins, "default"))
    setEdges(toEdges(sharedInfra.plugins, false))
  }, [])

  useEffect(() => {
    if (!name)
      return

    getSharedInfra(name)
  }, [name])

  const createSharedInfra = useCallback(async (sharedInfra: any) => {
    const res = await fetch(`http://localhost:8080/shared-infras`, {
      method: 'POST',
      body: JSON.stringify({
        ...sharedInfra,
        plugins,
      })
    })
    const created = await res.json()
    navigate(`/shared-infras/${sharedInfra?.name}`)
  }, [plugins])

  const handleDiagramChanges = useCallback((nodes: any, edges: any) => {
    let dict: any = {}
    for(let i = 0; i < nodes.length; i++) {
      dict[nodes[i].id] = nodes[i]?.data?.name
    }

    const newPlugins = nodes?.map((node: any) => {
      return {
        name: node?.data?.name,
        ref: node?.data?.ref,
        type: node?.data?.type,
        depends: edges.filter((e: any) => e.target === node.id).map((e: any) => dict[e.source]),
        inputs: node?.data?.inputs,
        outputs: [],
      }
    })

    setPlugins(newPlugins)
  }, [setPlugins])

  useEffect(() => {
    handleDiagramChanges(nodes, edges)
  }, [nodes, edges])

  useEffect(() => {
    console.log(plugins)
  }, [plugins])

  
  return (
    <div className="shared-infra-create__content">
      <DefaultPanel
        sharedInfra={sharedInfra}
        onSave={createSharedInfra}
        goToView={() => navigate(`/shared-infras/${name}`)}
      />
      <div className="shared-infra-view__diagram">
      <SharedInfraDiagram
        action="CREATE"
        nodes={nodes}
        edges={edges}
        onChangeDiagram={handleDiagramChanges}
      />
      </div>
    </div>
  )
})

export default SharedInfraEditor