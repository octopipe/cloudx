import React, { useCallback, useEffect, useState } from 'react';
import ReactFlow, { useNodesState, useEdgesState, addEdge, ConnectionLineType } from 'reactflow';
import ExecutionNode from './ExecutionNode';
import dagre from 'dagre'
import 'reactflow/dist/style.css';
import { Alert, Badge, Button, ListGroup, Modal } from 'react-bootstrap';

const dagreGraph = new dagre.graphlib.Graph();
dagreGraph.setDefaultEdgeLabel(() => ({}));

const nodeWidth = 172;
const nodeHeight = 70;

const getLayoutedElements = (nodes: any, edges: any) => {
  dagreGraph.setGraph({ rankdir: 'LR' });

  nodes.forEach((node: any) => {
    dagreGraph.setNode(node.id, { width: nodeWidth, height: nodeHeight });
  });

  edges.forEach((edge: any) => {
    dagreGraph.setEdge(edge.source, edge.target);
  });

  dagre.layout(dagreGraph);

  nodes.forEach((node: any) => {
    const nodeWithPosition = dagreGraph.node(node.id);
    node.targetPosition = 'left'
    node.sourcePosition = 'right'

    node.position = {
      x: nodeWithPosition.x - nodeWidth / 2,
      y: nodeWithPosition.y - nodeHeight / 2,
    };

    return node;
  });

  return { nodes, edges };
};

const nodeTypes = {
  executionNode: ExecutionNode,
};

const SharedInfraViewDiagram = ({ initialNodes, initialEdges }: any) => {
  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  const [selectedNode, setSelectedNode] = useState<any>()
  const [show, setShow] = useState(false);

  const handleClose = () => setShow(false);
  const handleShow = () => setShow(true);
  const onConnect = useCallback(
    (params: any) =>
      setEdges((eds) =>
        addEdge({ ...params, type: ConnectionLineType.SmoothStep, animated: true }, eds)
      ),
    []
  );
  const onNodeClick = (event: any, node: any) => {
    setSelectedNode(node)
    handleShow()
  }


  useEffect(() => {
    const { nodes: layoutedNodes, edges: layoutedEdges } = getLayoutedElements(
      initialNodes,
      initialEdges,
    );

    setNodes([...layoutedNodes]);
    setEdges([...layoutedEdges]);

  }, [initialNodes, initialEdges])
  
  return (
    <div style={{ width: '100%', height: '300px' }}>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        nodeTypes={nodeTypes}
        onNodeClick={onNodeClick}
        fitView
      ></ReactFlow>
      <Modal show={show} onHide={handleClose}>
        <Modal.Header closeButton>
          <Modal.Title>{selectedNode?.data?.label}</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          {selectedNode?.data?.error && (<Alert variant='danger'>
            {selectedNode?.data?.error?.replace(/(?:\\n|\\\\n)/g, '\n')}
          </Alert>)}
          <h1 className='h5'>Info</h1>
          <div className='mb-3'>
            <strong>Type: </strong><Badge>{selectedNode?.data?.type}</Badge><br/>
            <strong>Ref: </strong>{selectedNode?.data?.ref}
          </div>
          <h1 className='h5'>Inputs</h1>
          <ListGroup>
            { selectedNode?.data?.inputs?.map((i: any) => (
              <ListGroup.Item>
                <strong>{i?.key}: </strong>{i?.value}
              </ListGroup.Item>
            )) }
          </ListGroup>
          {/* <pre>{JSON.stringify(selectedNode?.data, null, 2)}</pre> */}
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={handleClose}>
            Close
          </Button>
          <Button variant="primary" onClick={handleClose}>
            Save Changes
          </Button>
        </Modal.Footer>
      </Modal>
    </div>
  )
}

export default SharedInfraViewDiagram