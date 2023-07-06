import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import Main from './modules/Main'
import SharedInfras from './modules/SharedInfras'
import ProvidersConfig from './modules/ProviderConfig'
import SharedInfraView from './modules/SharedInfrasView'
import ConnectionInterfaces from './modules/ConnectionInterfaces'
import './core/components/icon'
import SharedInfraEditor from './modules/SharedInfrasEditor'
import ProviderConfigView from './modules/ProviderConfigView'

const App = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Main />}>
          <Route path='' element={<Navigate to="/shared-infras" replace />} />
          <Route path='/shared-infras' element={<SharedInfras />} />
          <Route path='/shared-infras/create' element={<SharedInfraEditor />}/>
          <Route path='/shared-infras/:name/edit' element={<SharedInfraEditor />}/>
          <Route path='/shared-infras/:name' element={<SharedInfraView />}/>
          <Route path='/connection-interfaces' element={<ConnectionInterfaces />}/>
          <Route path='/providers-config' element={<ProvidersConfig />}/>
          <Route path='/providers-config/:name' element={<ProviderConfigView />}/>
        </Route>
        <Route path="login" element={<div>Login</div>} />
        <Route path="*" element={<div>Not Found</div>} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
