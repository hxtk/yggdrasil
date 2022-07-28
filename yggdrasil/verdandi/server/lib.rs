use std::io::Result;
use std::net::UdpSocket;

pub struct Server {
    socket: UdpSocket,
}

impl Server {
    fn serve(&self, addr: &str) -> Result<()> {
        let socket = UdpSocket::bind(addr)?;
        Ok(())
    }
}
