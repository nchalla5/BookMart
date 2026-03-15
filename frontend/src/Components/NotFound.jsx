function NotFound() {
    return (
        <div style={{minHeight: '100vh', display: 'grid', placeItems: 'center', padding: '24px'}}>
            <div style={{textAlign: 'center', background: 'rgba(255,255,255,0.82)', padding: '32px', borderRadius: '24px', boxShadow: '0 20px 80px rgba(31, 41, 55, 0.12)'}}>
                <p style={{letterSpacing: '0.16em', textTransform: 'uppercase', color: '#8c5f3c'}}>Book Mart</p>
                <h1>404! Page not found</h1>
            </div>
        </div>
    );
}
export default NotFound
