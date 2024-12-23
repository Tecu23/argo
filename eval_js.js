pos = {
  // chessboard
  b: [
    ["r", "p", "-", "-", "-", "-", "P", "R"],
    ["-", "p", "b", "-", "-", "-", "P", "-"],
    ["b", "p", "n", "-", "-", "N", "-", "B"],
    ["q", "-", "-", "p", "P", "B", "-", "Q"],
    ["-", "-", "p", "P", "-", "-", "-", "-"],
    ["r", "p", "-", "-", "-", "N", "P", "R"],
    ["k", "p", "-", "-", "-", "-", "P", "K"],
    ["-", "p", "-", "-", "-", "-", "P", "-"],
  ],

  // castling rights
  c: [false, false, false, false],

  // enpassant
  e: null,

  // side to move
  w: true,

  // move counts
  m: [0, 1],
};
function board(pos, x, y) {
  if (x >= 0 && x <= 7 && y >= 0 && y <= 7) return pos.b[x][y];
  return "x";
}
function colorflip(pos) {
  var board = new Array(8);
  for (var i = 0; i < 8; i++) board[i] = new Array(8);
  for (x = 0; x < 8; x++)
    for (y = 0; y < 8; y++) {
      board[x][y] = pos.b[x][7 - y];
      var color = board[x][y].toUpperCase() == board[x][y];
      board[x][y] = color
        ? board[x][y].toLowerCase()
        : board[x][y].toUpperCase();
    }
  return {
    b: board,
    c: [pos.c[2], pos.c[3], pos.c[0], pos.c[1]],
    e: pos.e == null ? null : [pos.e[0], 7 - pos.e[1]],
    w: !pos.w,
    m: [pos.m[0], pos.m[1]],
  };
}
function sum(pos, func, param) {
  var sum = 0;
  for (var x = 0; x < 8; x++)
    for (var y = 0; y < 8; y++) sum += func(pos, { x: x, y: y }, param);
  return sum;
}

function bishop_on_king_ring(pos, square) {
  if (square == null) return sum(pos, bishop_on_king_ring);
  if (board(pos, square.x, square.y) != "B") return 0;
  console.log("king attck", king_attackers_count(pos, square));
  if (king_attackers_count(pos, square) > 0) return 0;

  for (var i = 0; i < 4; i++) {
    var ix = (i > 1) * 2 - 1;
    var iy = (i % 2 == 0) * 2 - 1;
    for (var d = 1; d < 8; d++) {
      var x = square.x + d * ix,
        y = square.y + d * iy;
      if (board(pos, x, y) == "x") break;
      console.log(x, y, king_ring(pos, { x: x, y: y }));
      if (king_ring(pos, { x: x, y: y })) return 1;
      if (board(pos, x, y).toUpperCase() == "P") break;
    }
  }
  return 0;
}

function king_attackers_count(pos, square) {
  if (square == null) return sum(pos, king_attackers_count);
  if ("PNBRQ".indexOf(board(pos, square.x, square.y)) < 0) return 0;
  if (board(pos, square.x, square.y) == "P") {
    var v = 0;
    for (var dir = -1; dir <= 1; dir += 2) {
      var fr = board(pos, square.x + dir * 2, square.y) == "P";
      if (
        square.x + dir >= 0 &&
        square.x + dir <= 7 &&
        king_ring(pos, { x: square.x + dir, y: square.y - 1 }, true)
      )
        v = v + (fr ? 0.5 : 1);
    }

    console.log("First returned", v);
    return v;
  }
  for (var x = 0; x < 8; x++) {
    for (var y = 0; y < 8; y++) {
      var s2 = { x: x, y: y };

      if (king_ring(pos, s2)) {
        if (
          knight_attack(pos, s2, square) ||
          bishop_xray_attack(pos, s2, square) ||
          rook_xray_attack(pos, s2, square) ||
          queen_attack(pos, s2, square)
        ) {
          console.log(
            "Second returned",
            king_ring(pos, s2),
            knight_attack(pos, s2, square),
            bishop_xray_attack(pos, s2, square),
            rook_xray_attack(pos, s2, square),
            queen_attack(pos, s2, square),
            s2,
          );
          return 1;
        }
      }
    }
  }
  return 0;
}

function king_ring(pos, square, full) {
  if (square == null) return sum(pos, king_ring);
  if (
    !full &&
    board(pos, square.x + 1, square.y - 1) == "p" &&
    board(pos, square.x - 1, square.y - 1) == "p"
  )
    return 0;
  for (var ix = -2; ix <= 2; ix++) {
    for (var iy = -2; iy <= 2; iy++) {
      if (
        board(pos, square.x + ix, square.y + iy) == "k" &&
        ((ix >= -1 && ix <= 1) || square.x + ix == 0 || square.x + ix == 7) &&
        ((iy >= -1 && iy <= 1) || square.y + iy == 0 || square.y + iy == 7)
      )
        return 1;
    }
  }
  return 0;
}

function knight_attack(pos, square, s2) {
  if (square == null) return sum(pos, knight_attack);
  var v = 0;
  for (var i = 0; i < 8; i++) {
    var ix = ((i > 3) + 1) * ((i % 4 > 1) * 2 - 1);
    var iy = (2 - (i > 3)) * ((i % 2 == 0) * 2 - 1);
    var b = board(pos, square.x + ix, square.y + iy);
    if (
      b == "N" &&
      (s2 == null || (s2.x == square.x + ix && s2.y == square.y + iy)) &&
      !pinned(pos, { x: square.x + ix, y: square.y + iy })
    )
      v++;
  }
  return v;
}
function pinned(pos, square) {
  if (square == null) return sum(pos, pinned);
  if ("PNBRQK".indexOf(board(pos, square.x, square.y)) < 0) return 0;
  return pinned_direction(pos, square) > 0 ? 1 : 0;
}
function bishop_xray_attack(pos, square, s2) {
  if (square == null) return sum(pos, bishop_xray_attack);
  var v = 0;
  for (var i = 0; i < 4; i++) {
    var ix = (i > 1) * 2 - 1;
    var iy = (i % 2 == 0) * 2 - 1;
    for (var d = 1; d < 8; d++) {
      var b = board(pos, square.x + d * ix, square.y + d * iy);
      if (
        b == "B" &&
        (s2 == null || (s2.x == square.x + d * ix && s2.y == square.y + d * iy))
      ) {
        var dir = pinned_direction(pos, {
          x: square.x + d * ix,
          y: square.y + d * iy,
        });
        if (dir == 0 || Math.abs(ix + iy * 3) == dir) {
          console.log(
            dir,
            ix,
            iy,
            ix + iy * 3,
            d,
            square.x + d * ix,
            square.y + d * iy,
          );
          v++;
        }
      }
      if (b != "-" && b != "Q" && b != "q") break;
    }
  }
  console.log("Bishop score", v);
  return v;
}
function rook_xray_attack(pos, square, s2) {
  if (square == null) return sum(pos, rook_xray_attack);
  var v = 0;
  for (var i = 0; i < 4; i++) {
    var ix = i == 0 ? -1 : i == 1 ? 1 : 0;
    var iy = i == 2 ? -1 : i == 3 ? 1 : 0;
    for (var d = 1; d < 8; d++) {
      var b = board(pos, square.x + d * ix, square.y + d * iy);
      if (
        b == "R" &&
        (s2 == null || (s2.x == square.x + d * ix && s2.y == square.y + d * iy))
      ) {
        var dir = pinned_direction(pos, {
          x: square.x + d * ix,
          y: square.y + d * iy,
        });
        if (dir == 0 || Math.abs(ix + iy * 3) == dir) v++;
      }
      if (b != "-" && b != "R" && b != "Q" && b != "q") break;
    }
  }

  return v;
}

function queen_attack(pos, square, s2) {
  if (square == null) return sum(pos, queen_attack);
  var v = 0;
  for (var i = 0; i < 8; i++) {
    var ix = ((i + (i > 3)) % 3) - 1;
    var iy = (((i + (i > 3)) / 3) << 0) - 1;
    for (var d = 1; d < 8; d++) {
      var b = board(pos, square.x + d * ix, square.y + d * iy);
      if (
        b == "Q" &&
        (s2 == null || (s2.x == square.x + d * ix && s2.y == square.y + d * iy))
      ) {
        var dir = pinned_direction(pos, {
          x: square.x + d * ix,
          y: square.y + d * iy,
        });
        if (dir == 0 || Math.abs(ix + iy * 3) == dir) v++;
      }
      if (b != "-") break;
    }
  }
  return v;
}
function pinned_direction(pos, square) {
  if (square == null) return sum(pos, pinned_direction);
  if ("PNBRQK".indexOf(board(pos, square.x, square.y).toUpperCase()) < 0)
    return 0;
  var color = 1;
  if ("PNBRQK".indexOf(board(pos, square.x, square.y)) < 0) color = -1;
  for (var i = 0; i < 8; i++) {
    var ix = ((i + (i > 3)) % 3) - 1;
    var iy = (((i + (i > 3)) / 3) << 0) - 1;
    var king = false;
    for (var d = 1; d < 8; d++) {
      var b = board(pos, square.x + d * ix, square.y + d * iy);
      if (b == "K") king = true;
      if (b != "-") break;
    }
    if (king) {
      for (var d = 1; d < 8; d++) {
        var b = board(pos, square.x - d * ix, square.y - d * iy);
        if (
          b == "q" ||
          (b == "b" && ix * iy != 0) ||
          (b == "r" && ix * iy == 0)
        )
          return Math.abs(ix + iy * 3) * color;
        if (b != "-") break;
      }
    }
  }
  return 0;
}
console.log(bishop_on_king_ring(pos, { x: 3, y: 5 }));
