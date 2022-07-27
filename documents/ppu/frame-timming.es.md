
Lineas de Vertical Blank (241-260)
----------------------------------
El flag de VBlank de la PPU se activa en el segundo ciclo de la scanline 241.
En ese momento también ocurre NMI (si está activo).
La PPU no realiza accesos a memoria durante estas scanlines.
El flag VBlank se desactiva en el segundo ciclo de la scanline 261