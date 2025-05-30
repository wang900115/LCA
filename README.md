<!DOCTYPE html>
<html>
  <Body>
    <h1 align="center"><strong>Lightweight Chat Application</strong></h1> 
    <p align="center">
      <img src="assets/banner.png" alt="LCA Banner" width="300"/>
    </p>
    <section > 
      <h1>🚀 Overview</h1>
      <p>LCA is a simple and extensible WebSocket-based and RestFulAPI message chat system written in Go. It is designed to be easily deployable and modifiable for various real-time communication use cases.</p>
    </section>
    <section>
      <h1>🎯 Get Started </h1>
      <ol>
      <h2><li>Clone the repository</li></h2>
        <code>git clone https://github.com/wang900115/LCA.git</code>
      <h2><li>Change to the new path</li></h2>
        <code>cd LCA</code>
      <h2><li>Build application</li></h2>
        <ul>
        <h3><li>Linux</li></h3>
          <code>docker build -t lca . </code>
        <h3><li>Windows</li></h3>
          <code>go build ./cmd/LCA/main.go</code>
        </ul>
      <h2><li>Run application</li></h2>
        <ul>
        <h3><li>Linux</li></h3>
          <code>docker run --name lca-container -p 8080:8080 lca</code>
        <h3><li>Windows</li></h3>
          <code> ./main </code>
        </ul>
      </ol>
    </section>
    <section>
      <h1>✔ Licensing</h1>
      <p>Open Source License: Root and subdirectories.</p>
    </section>
  </Body>
</html>






