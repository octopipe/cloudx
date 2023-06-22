import React, { useCallback, useEffect } from 'react';
import ReactFlow, { useNodesState, useEdgesState, addEdge, ConnectionLineType } from 'reactflow';
import ExecutionNode from './ExecutionNode';
import dagre from 'dagre'
import 'reactflow/dist/style.css';

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
  const onConnect = useCallback(
    (params: any) =>
      setEdges((eds) =>
        addEdge({ ...params, type: ConnectionLineType.SmoothStep, animated: true }, eds)
      ),
    []
  );

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
        fitView
      ></ReactFlow>
    </div>
  )
}

export default SharedInfraViewDiagram