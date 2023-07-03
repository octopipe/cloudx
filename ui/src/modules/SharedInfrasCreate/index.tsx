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

const SharedInfraCreate = memo(() => {
  const navigate = useNavigate()
  const [nodes, setNodes] = useState([]);
  const [edges, setEdges] = useState([]);
  const [plugins, setPlugins] = useState<any>([])

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

  
  return (
    <div className="shared-infra-create__content">
      <DefaultPanel onCreate={createSharedInfra} />
      <SharedInfraDiagram
        action="CREATE"
        nodes={nodes}
        edges={edges}
        onChangeDiagram={handleDiagramChanges}
      />
    </div>
  )
})

export default SharedInfraCreate