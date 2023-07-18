import React, { memo, useEffect, useRef } from 'react';
import { Button } from 'react-bootstrap';
import { Handle, Position } from 'reactflow';
import h337, { BaseHeatmapConfiguration, HeatmapConfiguration } from 'heatmap.js'

const colorByStatus: any = {
  'SUCCESS': '#13aa80',
  'APPLIED': '#13aa80',
  'RUNNING': 'gray',
  'FAILED': 'red',
  'ERROR': 'red',
  'APPLY_ERROR': 'red'
}

const getDuration = (startedAt: any, finishedAt: any) => {
  const d1: any = new Date(startedAt)
  const d2: any = new Date(finishedAt)
  const diff = d2 - d1

  if (diff > 60e3) 
    return `${Math.floor(diff / 60e3)} minutes`

  return `${Math.floor(diff / 1e3)} seconds`
}

export default memo(({ data, isConnectable }: any) => {
  return (
    <>
      <Handle
        type="target"
        position={Position.Left}
        style={{ background: '#555' }}
        onConnect={(params) => console.log('handle onConnect', params)}
        isConnectable={isConnectable}
      />
      <div>
        <div 
          style={{
            position: 'absolute',
            width: '80px',
            height: '80px',
            zIndex: 5,
            background: "hsl(var(--hue), 100%, 50%"
          }}
        ></div>
        <div
          style={{
            background: data?.status ? colorByStatus[data?.status] : 'gray',
            color: '#fff',
            padding: '5px'
          }}
        >
          {data.label}
        </div>
        { data?.startedAt && data?.finishedAt && (
          <div style={{padding: '10px', color: "#000"}} id="execution-body" data-duration={getDuration(data?.startedAt, data?.finishedAt)}>

            {getDuration(data?.startedAt, data?.finishedAt)}
            {/* {(data?.status === "ERROR" || data?.status === "FAILED") && (
              <Button variant='danger'>See error</Button>
            )} */}

          </div>
        )}
        
      </div>
      <Handle
        type="source"
        position={Position.Right}
        id="a"
        style={{ background: '#555' }}
        isConnectable={isConnectable}
      />
      {
        data?.taskOutputs?.length > 0 && (
          <Handle
            type="source"
            position={Position.Right}
            id="cn"
            style={{ background: '#555', marginTop: '20px' }}
            isConnectable={isConnectable}
          />
        )
      }
      
    </>
  );
});