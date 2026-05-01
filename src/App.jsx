import React, { useState } from 'react';

function App() {
  const [activeTab, setActiveTab] = useState('MAPS');
  const [bannedMaps, setBannedMaps] = useState([]);
  const [isFinishing, setIsFinishing] = useState(false);

  const allMaps = ['Mirage', 'Inferno', 'Dust II', 'Ancient', 'Anubis', 'Vertigo', 'Nuke', 'Overpass', 'Train'];
  
  const handleBan = (map) => {
    if (bannedMaps.length < allMaps.length - 1) {
      const newBans = [...bannedMaps, map];
      setBannedMaps(newBans);
      
      if (newBans.length === allMaps.length - 1) {
        // Увеличиваем задержку, чтобы печать успела "упасть"
        setTimeout(() => setIsFinishing(true), 850); 
      }
    }
  };

  const winner = allMaps.find(m => !bannedMaps.includes(m));

  return (
    <div className="container">
      <div className="match-screen">
        <div className="teams-overview">
          <div className="team-box">
            <p className="t-name">TEAM A</p>
            {[1,2,3,4,5].map(i => <div key={i} className="p-item"><div className="p-avatar"></div>PLAYER_{i}</div>)}
          </div>
          <div className="team-box" style={{textAlign: 'right'}}>
            <p className="t-name">TEAM B</p>
            {[1,2,3,4,5].map(i => <div key={i} className="p-item" style={{flexDirection:'row-reverse'}}><div className="p-avatar" style={{marginLeft:'6px'}}></div>ENEMY_{i}</div>)}
          </div>
        </div>

        <div className="veto-grid">
          {allMaps.map((m, index) => {
            const isBanned = bannedMaps.includes(m);
            const isWinner = m === winner;
            const row = Math.floor(index / 3);
            const col = index % 3;
            
            let statusClass = "";
            if (isBanned) statusClass = "fade-out";
            if (isWinner && isFinishing) statusClass = "winner-move";

            return (
              <button 
                key={m} 
                style={{ left: `${col * 35}%`, top: `${row * 120}px` }}
                className={`map-btn ${statusClass}`}
                onClick={() => handleBan(m)}
              >
                <span>{m}</span>
                {isBanned && (
                  <div className="ban-overlay">
                    <div className="ban-stamp">BANNED</div>
                  </div>
                )}
              </button>
            );
          })}

          <div className={`server-reveal ${isFinishing ? 'visible' : ''}`}>
            <div className="srv-title">СЕРВЕР: БЕРЛИН</div>
            <p style={{color: '#8b949e', fontSize: '12px', marginBottom: '15px'}}>Матч готов. Скопируйте пароль лобби.</p>
            <div className="pass-block">
              <div className="pass-stars">****</div>
              <button className="copy-btn" onClick={() => alert('Пароль скопирован!')}>COPY</button>
            </div>
          </div>
        </div>
      </div>

      <nav className="bottom-nav">
        <button className="nav-item">🏠<br/>HOME</button>
        <button className="nav-item active">🗺️<br/>VETO</button>
      </nav>
    </div>
  );
}

export default App;
