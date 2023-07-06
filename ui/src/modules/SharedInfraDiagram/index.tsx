import React, { useCallback, useEffect, useRef, useState } from 'react';
import ReactFlow, { useNodesState, useEdgesState, addEdge, ConnectionLineType, Background, ReactFlowProvider, useReactFlow } from 'reactflow';
import DefaultNode from './DefaultNode'
import ExecutionNode from './ExecutionNode';
import dagre from 'dagre'
import 'reactflow/dist/style.css';
import './index.css'
import NodePanel from './NodePanel';
import AddPanel from './AddPanel';
import ConnectionInterfaceNode from './ConnectionInterfaceNode';

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
  default: DefaultNode,
  aws: DefaultNode,
  'connection-interface': ConnectionInterfaceNode,
  executionNode: ExecutionNode,
};

let id = 0;
const getId = () => `dndnode_${id++}`;

const SharedInfraDiagram = ({ sharedInfra, nodes: initialNodes, edges: initialEdges, action, onChangeDiagram }: any) => {
  const reactFlowWrapper = useRef<any>(null);
  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  const [selectedNode, setSelectedNode] = useState<any>()
  const [reactFlowInstance, setReactFlowInstance] = useState<any>(null);

  const onConnect = useCallback((params: any) => setEdges((eds) => addEdge(params, eds)), []);

  const onNodeClick = (event: any, node: any) => {
    setSelectedNode(node)
  }

  const onChangeNodePanel = (node: any) => {
    setNodes((nodes: any) => nodes.map((n: any) => {
      if (n.id == node.id) {
        return node
      }

      return n
    }))
  }

  useEffect(() => {
    const { nodes: layoutedNodes, edges: layoutedEdges } = getLayoutedElements(
      initialNodes,
      initialEdges,
    );

    setNodes(layoutedNodes);
    setEdges(layoutedEdges);
  }, [initialNodes, initialEdges])

  useEffect(() => {
    onChangeDiagram && onChangeDiagram(nodes, edges)
  }, [nodes, edges])

  const onDragOver = useCallback((event: any) => {
    event.preventDefault();
    event.dataTransfer.dropEffect = 'move';
  }, []);

  const onDrop = useCallback(
    (event: any) => {
      event.preventDefault();

      const reactFlowBounds = reactFlowWrapper.current.getBoundingClientRect();
      const type = event.dataTransfer.getData('application/reactflow');

      if (typeof type === 'undefined' || !type) {
        return;
      }

      const position = reactFlowInstance?.project({
        x: event.clientX - reactFlowBounds.left,
        y: event.clientY - reactFlowBounds.top,
      });
      const id = getId()
      const newNode = {
        id: id,
        type,
        position,
        data: { 
          label: id, 
          name: id, 
          category: type,
          ref: '',
          type: '',
          inputs: [
            { key: 'example-key', value: 'example-value' }
          ]
        },
      };

      setNodes((nds) => nds.concat(newNode));
    },
    [reactFlowInstance]
  );

  return (
    <div className='shared-infra-diagram'>
      <ReactFlowProvider>
        <div className="reactflow-wrapper" ref={reactFlowWrapper}>
          <ReactFlow
            nodes={nodes}
            edges={edges}
            nodeTypes={nodeTypes}
            onNodesChange={(e) => (action == "CREATE" || action == "UPDATE" ) ? onNodesChange(e) : null}
            onEdgesChange={(e) => (action == "CREATE" || action == "UPDATE" ) ? onEdgesChange(e) : null}
            onNodeClick={onNodeClick}
            onInit={setReactFlowInstance}
            onDrop={onDrop}
            onDragOver={onDragOver}
            onConnect={(e) => (action == "CREATE" || action == "UPDATE" ) ? onConnect(e) : null}
            fitView
          >
            <Background />
          </ReactFlow>
        </div>
      </ReactFlowProvider>
      { !!selectedNode && <NodePanel selectedNode={selectedNode} action={action} onClose={() => setSelectedNode(null)} onChange={onChangeNodePanel} /> }
      { (action == "CREATE" || action == "UPDATE" ) && !selectedNode && <AddPanel onClose={() => setSelectedNode(null)} /> }
    </div>
  )

}

export default (props: any) => (
  <ReactFlowProvider>
    <SharedInfraDiagram {...props}/>
  </ReactFlowProvider>
)