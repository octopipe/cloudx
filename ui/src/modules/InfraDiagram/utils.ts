const position = { x: 0, y: 0 };
const edgeType = 'smoothstep';


export const toNodes = (tasks: any, type="executionNode") => {
  return tasks.map((p: any) => {
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

}

export const toEdges = (tasks: any, animated: boolean) => {
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


  return edges
}