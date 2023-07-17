const position = { x: 0, y: 0 };
const edgeType = 'smoothstep';


export const toNodes = (tasks: any, type="executionNode") => {
  let nodes = tasks.map((p: any) => {
    return {
      id: p.name,
      type: type,
      sourcePosition: 'right',
      targetPosition: 'left',
      data: { 
        ...p,
        label: p.name,
        category: type,
      },
      position,
    }
  })

  if (type === "executionNode") {
    for(let i = 0; i < tasks.length; i++) {
      if (tasks[i]?.taskOutputs) {
        const taskOutputs = tasks[i]?.taskOutputs?.map((t: any) => ({
          id: t?.name,
          type: 'taskOutput',
          targetPosition: 'left',
          data: {
            label: t?.name,
          },
          position
        }))
  
        console.log(taskOutputs)
  
        nodes = [...nodes, ...taskOutputs]
      }
      
    }
  }
  return nodes
}

export const toEdges = (tasks: any, animated: boolean, type="executionNode") => {
  let edges: any = []
  for (let i = 0; i < tasks.length; i++) {
    for (let j = 0; j < tasks[i]?.depends?.length; j++) {
      edges = [...edges,  {
        id: `e-${tasks[i].name}-${tasks[i].depends[j]}`,
        source: tasks[i].depends[j],
        target: tasks[i].name,
        type: edgeType,
        animated,
      }]
    }
  }

  if (type === "executionNode") {
    for(let i = 0; i < tasks.length; i++) {
      if (tasks[i]?.taskOutputs) {
        const taskOutputs = tasks[i]?.taskOutputs?.map((t: any) => ({
          id: `e-${tasks[i].name}-${t.name}`,
          source: tasks[i].name,
          target: t.name,
          type: "default",
          sourceHandle: 'cn'
        }))
  
        edges = [...edges, ...taskOutputs]
      }
      
    }
  }


  return edges
}