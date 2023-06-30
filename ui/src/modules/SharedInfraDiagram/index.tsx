import React, { useCallback, useEffect, useRef, useState } from 'react';
import ReactFlow, { useNodesState, useEdgesState, addEdge, ConnectionLineType, Background } from 'reactflow';
import DefaultNode from './DefaultNode'
import ExecutionNode from './ExecutionNode';
import dagre from 'dagre'
import 'reactflow/dist/style.css';
import './index.css'
import NodePanel from './NodePanel';
import AddPanel from './AddPanel';

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
  defaultNode: DefaultNode,
  executionNode: ExecutionNode,
};

let id = 0;
const getId = () => `dndnode_${id++}`;

const SharedInfraDiagram = ({ nodes: initialNodes, edges: initialEdges, action }: any) => {
  const reactFlowWrapper = useRef<any>(null);
  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  const [selectedNode, setSelectedNode] = useState<any>()
  const [reactFlowInstance, setReactFlowInstance] = useState<any>(null);

  const onConnect = useCallback(
    (params: any) =>
      setEdges((eds) =>
        addEdge({ ...params, type: ConnectionLineType.SmoothStep, animated: true }, eds)
      ),
    []
  );
  const onNodeClick = (event: any, node: any) => {
    setSelectedNode(node)
  }
  const onClosePlane = () => {
    setSelectedNode(null)
  }


  useEffect(() => {
    const { nodes: layoutedNodes, edges: layoutedEdges } = getLayoutedElements(
      initialNodes,
      initialEdges,
    );

    setNodes([...layoutedNodes]);
    setEdges([...layoutedEdges]);

  }, [initialNodes, initialEdges])

  const onDragOver = useCallback((event: any) => {
    event.preventDefault();
    event.dataTransfer.dropEffect = 'move';
  }, []);

  const onDrop = useCallback(
    (event: any) => {
      event.preventDefault();

      const reactFlowBounds = reactFlowWrapper.current.getBoundingClientRect();
      const type = event.dataTransfer.getData('application/reactflow');

      // check if the dropped element is valid
      if (typeof type === 'undefined' || !type) {
        return;
      }

      const position = reactFlowInstance?.project({
        x: event.clientX - reactFlowBounds.left,
        y: event.clientY - reactFlowBounds.top,
      });
      const newNode = {
        id: getId(),
        type,
        position,
        data: { label: `${type} node` },
      };

      setNodes((nds) => nds.concat(newNode));
    },
    [reactFlowInstance]
  );

  return (
    <>
      <div className="reactflow-wrapper" ref={reactFlowWrapper} style={{width: '100%', height: '100%'}}>
        <ReactFlow
          nodes={nodes}
          edges={edges}
          nodeTypes={nodeTypes}
          onNodeClick={onNodeClick}
          onInit={setReactFlowInstance}
          onDrop={onDrop}
          onDragOver={onDragOver}
          fitView
        >
          <Background />
        </ReactFlow>
      </div>
      { !!selectedNode && <NodePanel node={selectedNode} onClose={() => setSelectedNode(null)} /> }
      { (action == "CREATE" || action == "UPDATE") && <AddPanel onClose={() => setSelectedNode(null)} /> }
    </>
  )

}

export default SharedInfraDiagram