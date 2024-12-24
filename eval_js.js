pos = {
  // chessboard
  b: [
    ["-", "p", "-", "-", "-", "-", "P", "-"],
    ["-", "p", "-", "-", "-", "-", "Q", "R"],
    ["r", "-", "n", "-", "-", "P", "-", "R"],
    ["r", "-", "-", "-", "P", "-", "-", "-"],
    ["-", "q", "p", "-", "-", "P", "P", "-"],
    ["-", "p", "n", "-", "-", "N", "P", "-"],
    ["k", "p", "-", "-", "-", "P", "B", "K"],
    ["-", "p", "-", "-", "-", "-", "P", "-"],
  ],
  // castling rights
  c: [false, false, false, false],

  // enpassant
  e: null,

  // side to move
  w: true,

  // move counts
  m: [0, 17],
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

function pieces_mg(pos, square) {
  if (square == null) return sum(pos, pieces_mg);
  if ("NBRQ".indexOf(board(pos, square.x, square.y)) < 0) return 0;

  var type = "NBRQ".indexOf(board(pos, square.x, square.y));
  var v = 0;
  v += 18 * minor_behind_pawn(pos, square);
  v -= 3 * bishop_pawns(pos, square);
  v -= 4 * bishop_xray_pawns(pos, square);
  v += 6 * rook_on_queen_file(pos, square);
  v += 16 * rook_on_king_ring(pos, square);
  v += 24 * bishop_on_king_ring(pos, square);
  v += [0, 19, 48][rook_on_file(pos, square)];
  v -= trapped_rook(pos, square) * 55 * (pos.c[0] || pos.c[1] ? 1 : 2);
  v -= 56 * weak_queen(pos, square);
  v -= 2 * queen_infiltration(pos, square);
  v -=
    (board(pos, square.x, square.y) == "N" ? 8 : 6) *
    king_protector(pos, square);
  v += 45 * long_diagonal_bishop(pos, square);
  return v;
}

function minor_behind_pawn(pos, square) {
  if (square == null) return sum(pos, minor_behind_pawn);
  if (
    board(pos, square.x, square.y) != "B" &&
    board(pos, square.x, square.y) != "N"
  )
    return 0;
  if (board(pos, square.x, square.y - 1).toUpperCase() != "P") return 0;
  return 1;
}

function bishop_pawns(pos, square) {
  if (square == null) return sum(pos, bishop_pawns);
  if (board(pos, square.x, square.y) != "B") return 0;
  var c = (square.x + square.y) % 2,
    v = 0;
  var blocked = 0;
  for (var x = 0; x < 8; x++) {
    for (var y = 0; y < 8; y++) {
      if (board(pos, x, y) == "P" && c == (x + y) % 2) v++;
      if (
        board(pos, x, y) == "P" &&
        x > 1 &&
        x < 6 &&
        board(pos, x, y - 1) != "-"
      )
        blocked++;
    }
  }
  return v * (blocked + (pawn_attack(pos, square) > 0 ? 0 : 1));
}

function pawn_attack(pos, square) {
  if (square == null) return sum(pos, pawn_attack);
  var v = 0;
  if (board(pos, square.x - 1, square.y + 1) == "P") v++;
  if (board(pos, square.x + 1, square.y + 1) == "P") v++;
  return v;
}

function bishop_xray_pawns(pos, square) {
  if (square == null) return sum(pos, bishop_xray_pawns);
  if (board(pos, square.x, square.y) != "B") return 0;
  var count = 0;
  for (var x = 0; x < 8; x++) {
    for (var y = 0; y < 8; y++) {
      if (
        board(pos, x, y) == "p" &&
        Math.abs(square.x - x) == Math.abs(square.y - y)
      )
        count++;
    }
  }
  return count;
}
function rook_on_queen_file(pos, square) {
  if (square == null) return sum(pos, rook_on_queen_file);
  if (board(pos, square.x, square.y) != "R") return 0;
  for (var y = 0; y < 8; y++) {
    if (board(pos, square.x, y).toUpperCase() == "Q") return 1;
  }
  return 0;
}
function rook_on_king_ring(pos, square) {
  if (square == null) return sum(pos, rook_on_king_ring);
  if (board(pos, square.x, square.y) != "R") {
    return 0;
  }
  if (king_attackers_count(pos, square) > 0) {
    return 0;
  }
  for (var y = 0; y < 8; y++) {
    if (king_ring(pos, { x: square.x, y: y })) return 1;
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
        )
          return 1;
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
  ) {
    return 0;
  }
  for (var ix = -2; ix <= 2; ix++) {
    for (var iy = -2; iy <= 2; iy++) {
      if (
        board(pos, square.x + ix, square.y + iy) == "k" &&
        ((ix >= -1 && ix <= 1) || square.x + ix == 0 || square.x + ix == 7) &&
        ((iy >= -1 && iy <= 1) || square.y + iy == 0 || square.y + iy == 7)
      ) {
        return 1;
      }
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
        if (dir == 0 || Math.abs(ix + iy * 3) == dir) v++;
      }
      if (b != "-" && b != "Q" && b != "q") break;
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

function bishop_on_king_ring(pos, square) {
  if (square == null) return sum(pos, bishop_on_king_ring);
  if (board(pos, square.x, square.y) != "B") return 0;
  if (king_attackers_count(pos, square) > 0) return 0;
  for (var i = 0; i < 4; i++) {
    var ix = (i > 1) * 2 - 1;
    var iy = (i % 2 == 0) * 2 - 1;
    for (var d = 1; d < 8; d++) {
      var x = square.x + d * ix,
        y = square.y + d * iy;
      if (board(pos, x, y) == "x") break;
      if (king_ring(pos, { x: x, y: y })) return 1;
      if (board(pos, x, y).toUpperCase() == "P") break;
    }
  }
  return 0;
}

function rook_on_file(pos, square) {
  if (square == null) return sum(pos, rook_on_file);
  if (board(pos, square.x, square.y) != "R") return 0;
  var open = 1;
  for (var y = 0; y < 8; y++) {
    if (board(pos, square.x, y) == "P") return 0;
    if (board(pos, square.x, y) == "p") open = 0;
  }
  return open + 1;
}

function trapped_rook(pos, square) {
  if (square == null) return sum(pos, trapped_rook);
  if (board(pos, square.x, square.y) != "R") return 0;
  if (rook_on_file(pos, square)) return 0;
  if (mobility(pos, square) > 3) return 0;
  var kx = 0,
    ky = 0;
  for (var x = 0; x < 8; x++) {
    for (var y = 0; y < 8; y++) {
      if (board(pos, x, y) == "K") {
        kx = x;
        ky = y;
      }
    }
  }
  if (kx < 4 != square.x < kx) return 0;
  return 1;
}

function weak_queen(pos, square) {
  if (square == null) return sum(pos, weak_queen);
  if (board(pos, square.x, square.y) != "Q") return 0;
  for (var i = 0; i < 8; i++) {
    var ix = ((i + (i > 3)) % 3) - 1;
    var iy = (((i + (i > 3)) / 3) << 0) - 1;
    var count = 0;
    for (var d = 1; d < 8; d++) {
      var b = board(pos, square.x + d * ix, square.y + d * iy);
      if (b == "r" && (ix == 0 || iy == 0) && count == 1) return 1;
      if (b == "b" && ix != 0 && iy != 0 && count == 1) return 1;
      if (b != "-") count++;
    }
  }
  return 0;
}

function queen_infiltration(pos, square) {
  if (square == null) return sum(pos, queen_infiltration);
  if (board(pos, square.x, square.y) != "Q") return 0;
  if (square.y > 3) return 0;
  if (board(pos, square.x + 1, square.y - 1) == "p") return 0;
  if (board(pos, square.x - 1, square.y - 1) == "p") return 0;
  if (pawn_attacks_span(pos, square)) return 0;
  return 1;
}

function pawn_attacks_span(pos, square) {
  if (square == null) return sum(pos, pawn_attacks_span);
  var pos2 = colorflip(pos);
  for (var y = 0; y < square.y; y++) {
    if (
      board(pos, square.x - 1, y) == "p" &&
      (y == square.y - 1 ||
        (board(pos, square.x - 1, y + 1) != "P" &&
          !backward(pos2, { x: square.x - 1, y: 7 - y })))
    )
      return 1;
    if (
      board(pos, square.x + 1, y) == "p" &&
      (y == square.y - 1 ||
        (board(pos, square.x + 1, y + 1) != "P" &&
          !backward(pos2, { x: square.x + 1, y: 7 - y })))
    )
      return 1;
  }
  return 0;
}

function king_protector(pos, square) {
  if (square == null) return sum(pos, king_protector);
  if (
    board(pos, square.x, square.y) != "N" &&
    board(pos, square.x, square.y) != "B"
  )
    return 0;
  return king_distance(pos, square);
}

function king_distance(pos, square) {
  if (square == null) return sum(pos, king_distance);
  for (var x = 0; x < 8; x++) {
    for (var y = 0; y < 8; y++) {
      if (board(pos, x, y) == "K") {
        return Math.max(Math.abs(x - square.x), Math.abs(y - square.y));
      }
    }
  }
  return 0;
}

function long_diagonal_bishop(pos, square) {
  if (square == null) return sum(pos, long_diagonal_bishop);
  if (board(pos, square.x, square.y) != "B") return 0;
  if (square.x - square.y != 0 && square.x - (7 - square.y) != 0) return 0;
  var x1 = square.x,
    y1 = square.y;
  if (Math.min(x1, 7 - x1) > 2) return 0;
  for (var i = Math.min(x1, 7 - x1); i < 4; i++) {
    if (board(pos, x1, y1) == "p") return 0;
    if (board(pos, x1, y1) == "P") return 0;
    if (x1 < 4) x1++;
    else x1--;
    if (y1 < 4) y1++;
    else y1--;
  }
  return 1;
}

function mobility(pos, square) {
  if (square == null) return sum(pos, mobility);
  var v = 0;
  var b = board(pos, square.x, square.y);
  if ("NBRQ".indexOf(b) < 0) return 0;
  for (var x = 0; x < 8; x++) {
    for (var y = 0; y < 8; y++) {
      var s2 = { x: x, y: y };
      if (!mobility_area(pos, s2)) continue;
      if (b == "N" && knight_attack(pos, s2, square) && board(pos, x, y) != "Q")
        v++;
      if (
        b == "B" &&
        bishop_xray_attack(pos, s2, square) &&
        board(pos, x, y) != "Q"
      )
        v++;
      if (b == "R" && rook_xray_attack(pos, s2, square)) v++;
      if (b == "Q" && queen_attack(pos, s2, square)) v++;
    }
  }
  return v;
}
function mobility_area(pos, square) {
  if (square == null) return sum(pos, mobility_area);
  if (board(pos, square.x, square.y) == "K") return 0;
  if (board(pos, square.x, square.y) == "Q") return 0;
  if (board(pos, square.x - 1, square.y - 1) == "p") return 0;
  if (board(pos, square.x + 1, square.y - 1) == "p") return 0;
  if (
    board(pos, square.x, square.y) == "P" &&
    (rank(pos, square) < 4 || board(pos, square.x, square.y - 1) != "-")
  )
    return 0;
  if (blockers_for_king(colorflip(pos), { x: square.x, y: 7 - square.y }))
    return 0;
  return 1;
}
function rank(pos, square) {
  if (square == null) return sum(pos, rank);
  return 8 - square.y;
}

function file(pos, square) {
  if (square == null) return sum(pos, file);
  return 1 + square.x;
}

function non_pawn_material(pos, square) {
  if (square == null) return sum(pos, non_pawn_material);
  var i = "NBRQ".indexOf(board(pos, square.x, square.y));
  if (i >= 0) return piece_value_bonus(pos, square, true);
  return 0;
}

function piece_value_bonus(pos, square, mg) {
  if (square == null) return sum(pos, piece_value_bonus);
  var a = mg ? [124, 781, 825, 1276, 2538] : [206, 854, 915, 1380, 2682];
  var i = "PNBRQ".indexOf(board(pos, square.x, square.y));
  if (i >= 0) return a[i];
  return 0;
}

function blockers_for_king(pos, square) {
  if (square == null) return sum(pos, blockers_for_king);
  if (pinned_direction(colorflip(pos), { x: square.x, y: 7 - square.y }))
    return 1;
  return 0;
}
function king_attack(pos, square) {
  if (square == null) return sum(pos, king_attack);
  for (var i = 0; i < 8; i++) {
    var ix = ((i + (i > 3)) % 3) - 1;
    var iy = (((i + (i > 3)) / 3) << 0) - 1;
    if (board(pos, square.x + ix, square.y + iy) == "K") return 1;
  }

  return 0;
}
function attack(pos, square) {
  if (square == null) return sum(pos, attack);
  var v = 0;
  v += pawn_attack(pos, square);
  v += king_attack(pos, square);
  v += knight_attack(pos, square);
  v += bishop_xray_attack(pos, square);
  v += rook_xray_attack(pos, square);
  v += queen_attack(pos, square);
  return v;
}
function space(pos, square) {
  if (non_pawn_material(pos) + non_pawn_material(colorflip(pos)) < 12222)
    return 0;
  var pieceCount = 0,
    blockedCount = 0;
  for (var x = 0; x < 8; x++) {
    for (var y = 0; y < 8; y++) {
      if ("PNBRQK".indexOf(board(pos, x, y)) >= 0) pieceCount++;
      if (
        board(pos, x, y) == "P" &&
        (board(pos, x, y - 1) == "p" ||
          (board(pos, x - 1, y - 2) == "p" && board(pos, x + 1, y - 2) == "p"))
      )
        blockedCount++;
      if (
        board(pos, x, y) == "p" &&
        (board(pos, x, y + 1) == "P" ||
          (board(pos, x - 1, y + 2) == "P" && board(pos, x + 1, y + 2) == "P"))
      )
        blockedCount++;
    }
  }
  var weight = pieceCount - 3 + Math.min(blockedCount, 9);
  if (space_area(pos, square) > 0) {
    console.log(square, space_area(pos, square), pieceCount, blockedCount);
  }
  return ((space_area(pos, square) * weight * weight) / 16) << 0;
}

function space_area(pos, square) {
  if (square == null) return sum(pos, space_area);
  var v = 0;
  var r = rank(pos, square);
  var f = file(pos, square);
  if (
    r >= 2 &&
    r <= 4 &&
    f >= 3 &&
    f <= 6 &&
    board(pos, square.x, square.y) != "P" &&
    board(pos, square.x - 1, square.y - 1) != "p" &&
    board(pos, square.x + 1, square.y - 1) != "p"
  ) {
    v++;
    if (
      (board(pos, square.x, square.y - 1) == "P" ||
        board(pos, square.x, square.y - 2) == "P" ||
        board(pos, square.x, square.y - 3) == "P") &&
      !attack(colorflip(pos), { x: square.x, y: 7 - square.y })
    )
      v++;
  }
  return v;
}

console.log(space(pos));
