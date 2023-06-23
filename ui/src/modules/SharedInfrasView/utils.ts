const position = { x: 0, y: 0 };
const edgeType = 'smoothstep';


export const toNodes = (plugins: any, type="executionNode") => {
  return plugins.map((p: any) => {
    return {
      id: p.name,
      type: type,
      sourcePosition: 'right',
      targetPosition: 'left',
      data: { 
        label: p.name,
        status: p?.status,
        startedAt: p?.startedAt,
        finishedAt: p?.finishedAt,
        error: p?.error
      },
      position,
    }
  })

}

export const toEdges = (plugins: any) => {
  let edges: any = []
  for (let i = 0; i < plugins.length; i++) {
    for (let j = 0; j < plugins[i]?.depends?.length; i++) {
      edges = [...edges,  {
        id: `e-${plugins[i].name}-${plugins[i].depends[j]}`,
        source: plugins[i].depends[j],
        target: plugins[i].name,
        type: edgeType,
        animated: true
      }]
    }
  }


  return edges
}