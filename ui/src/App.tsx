import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import Main from './modules/Main'
import SharedInfras from './modules/SharedInfras'
import ProvidersConfig from './modules/ProviderConfig'
import SharedInfraView from './modules/SharedInfrasView'
import ConnectionInterfaces from './modules/ConnectionInterfaces'
import './core/components/icon'
import SharedInfraCreate from './modules/SharedInfrasCreate'

const App = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Main />}>
          <Route path='' element={<Navigate to="/shared-infras" replace />} />
          <Route path='/shared-infras' element={<SharedInfras />} />
          <Route path='/shared-infras/create' element={<SharedInfraCreate />}/>
          <Route path='/shared-infras/:name' element={<SharedInfraView />}/>
          <Route path='/connection-interfaces' element={<ConnectionInterfaces />}/>
          <Route path='/providers-config' element={<ProvidersConfig />}/>
        </Route>
        <Route path="login" element={<div>Login</div>} />
        <Route path="*" element={<div>Not Found</div>} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
