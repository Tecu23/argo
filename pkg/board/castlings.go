package board

// Castlings represents the castling possibilities for a given position
type Castlings uint

/*
                                castling  move    in       in
                                   right update   binary   decimal

   king & rooks didn't move:       1111 & 1111  =  1111    15

           white king moved:       1111 & 1100  =  1100    12
    white king's rook moved:       1111 & 1110  =  1110    14
   white queen's rook moved:       1111 & 1101  =  1101    13

           black king moved:       1111 & 0011  =  0011    3
    black king's rook moved:       1111 & 1011  =  1011    11
   black queen's rook moved:       1111 & 0111  =  0111    7
*/
